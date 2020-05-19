package storage

import (
	"github.com/CoverGenius/backup/base"
	h "github.com/CoverGenius/backup/helpers"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
	"strings"
	"sync"
	"time"
)

// YYYY.MM.DD.HH.mm.SS
const time_format = "2006.01.02.15.04.05"

type S3 struct {
	Session  *s3.S3
	Provider *session.Session
}

func (s *S3) Verify(c *base.Config) error {
	if c.Storage.S3.AWSMeta.AccessKey == nil || c.Storage.S3.AWSMeta.SecretAccessKey == nil {
		return errors.New("Access key and secret access key MUST be specified!")
	}
	if c.Storage.S3.Region == nil {
		return errors.New("Region MUST be specified!")
	}
	if c.Storage.S3.Bucket == nil {
		return errors.New("Bucket MUST be specified!")
	}

	return nil
}

func (s *S3) Pre(c *base.Config) error {
	s.Provider = session.Must(session.NewSession(&aws.Config{
		Region: c.Storage.S3.Region,
		Credentials: credentials.NewStaticCredentials(
			*c.Storage.S3.AWSMeta.AccessKey,
			*c.Storage.S3.AWSMeta.SecretAccessKey,
			"",
		),
	}))
	s.Session = s3.New(s.Provider)

	return nil
}

func (s *S3) Store(c *base.Config) error {
	h.Log.Debug("Store backup on S3 ...")

	s3_key := fmt.Sprintf(
		"%s/%s/backup.%s.%s.%s", *c.Name, time.Now().Format(time_format),
		*c.ArchiverSuffix, *c.CompressorSuffix, *c.EncryptorSuffix,
	)
	fPath := fmt.Sprintf("%s/backup.%s.%s.%s", *c.TmpDir, *c.ArchiverSuffix, *c.CompressorSuffix, *c.EncryptorSuffix)

	f, err := os.Open(fPath)
	defer f.Close()
	h.LogErrorExit(err)

	uploader := s3manager.NewUploader(s.Provider)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: c.Storage.S3.Bucket,
		Key:    &s3_key,
		Body:   f,
	})
	h.LogErrorExit(err)

	s.Session.WaitUntilObjectExists(&s3.HeadObjectInput{
		Bucket: c.Storage.S3.Bucket,
		Key:    &s3_key,
	})

	return nil
}

func FindObject(s *S3, c *base.Config) *string {
	// If key is not specified assume user wants
	// to restore from latest version
	if c.Storage.Key == nil {
		c.Storage.Key = h.StringP("latest")
	}

	result, err := s.Session.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: c.Storage.S3.Bucket,
		Prefix: c.Name,
	})
	h.LogError(err)

	if len(result.Contents) < 1 {
		h.LogErrorExit(fmt.Errorf("Found 0 backups for ID: %s!", *c.Name))
	}

	data := []*h.TData{}
	for _, item := range result.Contents {
		// If user explicitly defined backup key and it was found than return it
		if *c.Storage.Key != "latest" && strings.Contains(*item.Key, *c.Storage.Key) == true {
			return item.Key
		}
		d := h.TData{
			Timestamp: (*item.LastModified).Unix(),
			Data:      *item.Key,
		}
		data = append(data, &d)
	}
	// If user explicitly defined backup, at this stage we can assume it was not found
	if *c.Storage.Key != "latest" {
		h.LogErrorExit(fmt.Errorf("No backup found with key: %s!", *c.Name))
	}
	h.TDataQuickSort(data)

	return &data[0].Data
}

func (s *S3) Fetch(c *base.Config) error {
	h.Log.Debug("Fetch data from S3 ...")

	s3_key := FindObject(s, c)
	parts := strings.Split(*s3_key, "/")

	c.MorphFilename = h.GetLastStringElement(parts)
	fName := fmt.Sprintf("%s/%s", *c.TmpDir, *c.MorphFilename)
	f, err := os.Create(fName)
	defer f.Close()
	h.LogErrorExit(err)

	downloader := s3manager.NewDownloader(s.Provider)
	_, err = downloader.Download(f, &s3.GetObjectInput{
		Bucket: c.Storage.S3.Bucket,
		Key:    s3_key,
	})
	h.LogErrorExit(err)

	return nil
}

func DeleteObject(bucket *string, key *string, wg *sync.WaitGroup, s *s3.S3) {
	delete_object_input := &s3.DeleteObjectInput{
		Bucket: bucket,
		Key:    key,
	}
	describe_head_object_input := &s3.HeadObjectInput{
		Bucket: bucket,
		Key:    key,
	}
	h.Log.Debug("Found item with key: ", *key, " violating retention policy! Deleting it ...")
	_, err := s.DeleteObject(delete_object_input)
	h.LogError(err)

	err = s.WaitUntilObjectNotExists(describe_head_object_input)
	h.LogError(err)

	wg.Done()
}

func (s *S3) Post(c *base.Config) error {
	var wg sync.WaitGroup

	result, err := s.Session.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: c.Storage.S3.Bucket,
		Prefix: c.Name,
	})
	h.LogError(err)

	toDelete := []*h.TData{}
	if len(result.Contents) > int(*c.Keep) {
		for _, item := range result.Contents {
			d := h.TData{
				Timestamp: (*item.LastModified).Unix(),
				Data:      *item.Key,
			}
			toDelete = append(toDelete, &d)
		}
	}
	h.TDataQuickSort(toDelete)
	if len(toDelete) > int(*c.Keep) {
		for _, item := range toDelete[*c.Keep:] {
			wg.Add(1)
			go DeleteObject(c.Storage.S3.Bucket, &item.Data, &wg, s.Session)
		}
	}
	wg.Wait()

	return nil
}
