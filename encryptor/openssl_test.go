package encryptor

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

func TestEncryptDecryptOpenSSL(test *testing.T) {
	h.Log = logrus.New()
	encrypt := OpenSSL{}
	config := base.Config{}
	config.Name = h.StringP("test-backup")
	config.TmpDir = h.StringP("/tmp/backup-test")
	config.ArchiverSuffix = h.StringP("archiver")
	config.CompressorSuffix = h.StringP("compressor")
	config.EncryptorSuffix = h.StringP("enc")
	config.Encryptor.OpenSSL.Password = h.StringP("sup3rs3cr3t!")
	config.MorphFilename = h.StringP("backup.archiver.compressor.enc")
	os.Mkdir(*config.TmpDir, 0700)
	h.Log = logrus.New()

	fName := fmt.Sprintf("%s/backup.archiver.compressor", *config.TmpDir)
	encryptedFilename := fmt.Sprintf("%s/backup.archiver.compressor.enc", *config.TmpDir)

	f, _ := os.Create(fName)
	f.WriteString(testContent)
	f.Close()

	// Encrypt
	encrypt.Encrypt(&config)
	os.Remove(fName)

	if h.IsFileExists(&encryptedFilename) == false {
		test.Errorf("Failed to create an encrypted file: %s", encryptedFilename)
	}

	// Decrypt
	encrypt.Decrypt(&config)
	if h.IsFileExists(&fName) == false {
		test.Errorf("File: %s does not exists after decrypt operation!", fName)
	}
	content, _ := ioutil.ReadFile(fName)
	if string(content) != testContent {
		test.Errorf("Result is incorrect, got: %s, want: %s.", content, testContent)
	}

	h.RemoveDir(config.TmpDir)
}
