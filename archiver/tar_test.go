package archiver

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

func TestArchiveUnarchiveTAR(test *testing.T) {
	h.Log = logrus.New()
	archive := TAR{}
	config := base.Config{}
	config.Name = h.StringP("test-backup")
	config.TmpDir = h.StringP("/tmp/backup-test")
	config.ArchiverSuffix = h.StringP("tar")
	config.MorphFilename = h.StringP("backup.tar")
	os.Mkdir(*config.TmpDir, 0700)
	h.Log = logrus.New()

	fName := fmt.Sprintf("%s/data.txt", *config.TmpDir)
	archiveName := fmt.Sprintf("%s/backup.tar", *config.TmpDir)

	f, _ := os.Create(fName)
	f.WriteString(testContent)
	f.Close()

	// Archive
	archive.Archive(&config)
	os.Remove(fName)

	if h.IsFileExists(&archiveName) == false {
		test.Errorf("Failed to create an archive: %s", archiveName)
	}

	fi, _ := os.Stat(archiveName)
	if fi.Size() < 10000 {
		test.Errorf("Archive: %s is too small!", archiveName)
	}

	// Unarchive
	archive.Unarchive(&config)
	if h.IsFileExists(&fName) == false {
		test.Errorf("File: %s does not exists after unarchive operation!", fName)
	}
	content, _ := ioutil.ReadFile(fName)
	if string(content) != testContent {
		test.Errorf("Result is incorrect, got: %s, want: %s.", content, testContent)
	}

	h.RemoveDir(config.TmpDir)
}
