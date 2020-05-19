package compressor

import (
	"github.com/CoverGenius/backup/base"
	h "github.com/CoverGenius/backup/helpers"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"testing"
)

var (
	testContent string = "ABCD1234abc!"
)

func TestCompressDecompressGZIP(test *testing.T) {
	h.Log = logrus.New()
	compress := GZIP{}
	config := base.Config{}
	config.Name = h.StringP("test-backup")
	config.TmpDir = h.StringP("/tmp/backup-test")
	config.ArchiverSuffix = h.StringP("txt")
	config.MorphFilename = h.StringP("backup.txt.gz")
	os.Mkdir(*config.TmpDir, 0700)
	h.Log = logrus.New()

	fName := fmt.Sprintf("%s/backup.txt", *config.TmpDir)
	compressedFilename := fmt.Sprintf("%s/backup.txt.gz", *config.TmpDir)

	f, _ := os.Create(fName)
	f.WriteString(testContent)
	f.Close()

	// Compress
	compress.Compress(&config)
	os.Remove(fName)

	if h.IsFileExists(&compressedFilename) == false {
		test.Errorf("Failed to create compressed file: %s", compressedFilename)
	}

	fi, _ := os.Stat(compressedFilename)
	if fi.Size() > 64 {
		test.Errorf("Compressed file: %s is too big!", compressedFilename)
	}

	// Decompress
	compress.Decompress(&config)
	if h.IsFileExists(&fName) == false {
		test.Errorf("File: %s does not exists after decompress operation!", fName)
	}
	content, _ := ioutil.ReadFile(fName)
	if string(content) != testContent {
		test.Errorf("Result is incorrect, got: %s, want: %s.", content, testContent)
	}

	h.RemoveDir(config.TmpDir)
}
