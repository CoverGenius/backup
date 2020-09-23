package main

import (
	"github.com/CoverGenius/backup/archiver"
	"github.com/CoverGenius/backup/base"
	"github.com/CoverGenius/backup/compressor"
	"github.com/CoverGenius/backup/encryptor"
	h "github.com/CoverGenius/backup/helpers"
	"github.com/CoverGenius/backup/notifier"
	"github.com/CoverGenius/backup/source/database"
	"github.com/CoverGenius/backup/source/local"
	"github.com/CoverGenius/backup/splitter"
	"github.com/CoverGenius/backup/storage"
	"flag"
	"fmt"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

var (
	config string
)

func init() {
	flag.StringVar(&config, "config", "config.yaml", "Path to config file")
	h.Log = logrus.New()
}

func processConfig(p *string) *base.Config {
	c := base.Config{}
	h.YAMLDecode(p, &c)

	if c.TmpDir == nil {
		fmt.Println("No tmp directory defined! Using default one")
		c.TmpDir = h.StringP("/tmp/backup")
	}
        if c.TmpDirCleanup != nil && *c.TmpDirCleanup == false {
                h.Log.Info("Skip cleanup of temporary directory!")
        } else {
                h.RemoveDir(c.TmpDir)
        }
	os.Mkdir(*c.TmpDir, 0700)
	h.Config = &c

	return &c
}
func Backup(c *base.Config) {
	var backup base.Source
	t := strings.ToLower(*c.Backup.Type)
	if t == "rds" {
		backup = &database.RDS{}
	} else if t == "redis" {
		backup = &database.Redis{}
	} else if t == "mysql" {
		backup = &database.MySQL{}
	} else if t == "elasticsearch" {
		backup = &database.Elasticsearch{}
	} else if t == "postgres" {
		backup = &database.Postgres{}
	} else if t == "rabbitmq" {
		backup = &database.RabbitMQ{}
	} else if t == "directory" {
		backup = &local.Directory{}
	}
	err := backup.Verify(c)
	h.LogErrorExit(err)
	backup.Pre(c)

	if c.Restore != nil && *c.Restore == true {
		backup.Restore(c)
	} else {
		backup.Backup(c)
		backup.Post(c)
	}
}

func Archive(c *base.Config) {
	var archive base.Archiver
	if c.Archiver.Type != nil {
		t := strings.ToLower(*c.Archiver.Type)
		if t == "tar" {
			archive = &archiver.TAR{}
			c.ArchiverSuffix = h.StringP("tar")
		}
		err := archive.Verify(c)
		h.LogErrorExit(err)
		archive.Pre(c)

		if c.Restore != nil && *c.Restore == true {
			archive.Unarchive(c)
		} else {
			archive.Archive(c)
			archive.Post(c)
		}
	}
}

func Compress(c *base.Config) {
	var compress base.Compressor
	if c.Compressor.Type != nil {
		t := strings.ToLower(*c.Compressor.Type)
		if t == "gzip" {
			c.CompressorSuffix = h.StringP("gz")
			compress = &compressor.GZIP{}
		} else if t == "bzip2" {
			compress = &compressor.BZIP2{}
		} else if t == "pigz" {
			compress = &compressor.PIGZ{}
		}
		err := compress.Verify(c)
		h.LogErrorExit(err)
		compress.Pre(c)
		if c.Restore != nil && *c.Restore == true {
			compress.Decompress(c)
		} else {
			compress.Compress(c)
			compress.Post(c)
		}
	}
}

func Encrypt(c *base.Config) {
	var encrypt base.Encryptor
	if c.Encryptor.Type != nil {
		t := strings.ToLower(*c.Encryptor.Type)
		if t == "gpg" {
			encrypt = &encryptor.GPG{}
		} else if t == "openssl" {
			c.EncryptorSuffix = h.StringP("enc")
			encrypt = &encryptor.OpenSSL{}
		}
		err := encrypt.Verify(c)
		h.LogErrorExit(err)
		encrypt.Pre(c)
		if c.Restore != nil && *c.Restore == true {
			encrypt.Decrypt(c)
		} else {
			encrypt.Encrypt(c)
			encrypt.Post(c)
		}
	}
}

func Split(c *base.Config) {
	var split base.Splitter
	if c.Splitter.Type != nil {
		t := strings.ToLower(*c.Splitter.Type)
		if t == "default" {
			split = &splitter.Split{}
		}
		err := split.Verify(c)
		h.LogErrorExit(err)
		split.Pre(c)
		split.Split(c)
		split.Post(c)
	}
}

func Store(c *base.Config) {
	var store base.Storage
	if c.Storage.Type != nil {
		t := strings.ToLower(*c.Storage.Type)
		if t == "sftp" {
			store = &storage.SFTP{}
		} else if t == "s3" {
			store = &storage.S3{}
		} else if t == "path" {
			store = &storage.Local{}
		}
		err := store.Verify(c)
		h.LogErrorExit(err)
		store.Pre(c)
		if c.Restore != nil && *c.Restore == true {
			store.Fetch(c)
		} else {
			store.Store(c)
			store.Post(c)
		}
	}
}

func Notify(c *base.Config) {
	if c.Notifier.Slack.WebhookURL != nil {
		s := &notifier.Slack{}
		h.Notifiers = append(h.Notifiers, s)
	}
	if c.Notifier.Mail.To != nil {
		if c.Notifier.Mail.Address == nil {
			h.Log.Warn("No SMTP server specified. Using default: 127.0.0.1:25")
			c.Notifier.Mail.Address = h.StringP("127.0.0.1:25")
		}
		m := &notifier.Mail{}
		h.Notifiers = append(h.Notifiers, m)
	}
	for _, n := range h.Notifiers {
		n.Verify(c)
		n.Pre(c)
	}
}

func RunBackup(c *base.Config) {
	Notify(c)
	Backup(c)
	Archive(c)
	Compress(c)
	Encrypt(c)
	Split(c)
	Store(c)
}

func RunRestore(c *base.Config) {
	Notify(c)
	Store(c)
	Split(c)
	Encrypt(c)
	Compress(c)
	Archive(c)
	Backup(c)
}

func SetLogger(c *base.Config) error {
	h.Log.SetOutput(os.Stdout)
	level, err := logrus.ParseLevel(*c.LogLevel)
	h.LogErrorExit(err)
	h.Log.SetLevel(level)
	c.LogFile = h.StringP(fmt.Sprintf("%s/%s.log", *c.TmpDir, *c.Name))
	os.Remove(*c.LogFile)
	h.Log.Hooks.Add(lfshook.NewHook(
		*c.LogFile,
		&logrus.TextFormatter{},
	))

	return nil
}

func main() {
	flag.Parse()
	c := processConfig(&config)
	SetLogger(c)
	c.Notifier.Status.Status = h.StringP("OK")
	c.Notifier.Status.StartTime = h.GetTimeNow()
	if c.Restore != nil && *c.Restore == true {
		RunRestore(c)
	} else {
		RunBackup(c)
	}
	if c.TmpDirCleanup != nil && *c.TmpDirCleanup == false {
		h.Log.Info("Skip cleanup of temporary directory!")
	} else {
		h.RemoveDir(c.TmpDir)
	}
	c.Notifier.Status.EndTime = h.GetTimeNow()
	c.Notifier.Status.CalculateDuration()
	for _, n := range h.Notifiers {
		n.Notify(c)
		n.Post(c)
	}
}
