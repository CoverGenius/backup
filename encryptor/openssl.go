package encryptor

import (
	"github.com/CoverGenius/backup/base"
	h "github.com/CoverGenius/backup/helpers"
	"errors"
	"fmt"
	"os"
	"strings"
)

type OpenSSL struct {
}

func (o *OpenSSL) Verify(c *base.Config) error {
	if h.IsCommandExists(h.StringP("openssl")) == false {
		h.LogErrorExit(errors.New("Please install openssl!"))
	}
	if c.Encryptor.OpenSSL.Password == nil {
		h.LogErrorExit(errors.New("Please specify encryption password!"))
	}

	return nil
}

func (o *OpenSSL) Pre(c *base.Config) error {
	return nil
}

func (o *OpenSSL) Encrypt(c *base.Config) error {
	h.Log.Debug("Encrypting data using OpenSSL ...")

	in := fmt.Sprintf("%s/backup.%s.%s", *c.TmpDir, *c.ArchiverSuffix, *c.CompressorSuffix)
	out := fmt.Sprintf("%s/backup.%s.%s.%s", *c.TmpDir, *c.ArchiverSuffix, *c.CompressorSuffix, *c.EncryptorSuffix)

	args := []string{
		"aes-256-cbc", "-e", "-a",
		"-md", "sha256",
		"-pbkdf2", "-iter", "20000",
		"-k", *c.Encryptor.OpenSSL.Password,
		"-in", in,
		"-out", out,
	}
	h.RunCommand("openssl", args)

	return nil
}

func (o *OpenSSL) Decrypt(c *base.Config) error {
	h.Log.Debug("Decrypting data using OpenSSL ...")

	in := fmt.Sprintf("%s/%s", *c.TmpDir, *c.MorphFilename)
	parts := strings.Split(*c.MorphFilename, ".")

	if *(h.GetLastStringElement(parts)) != "enc" {
		h.LogErrorExit(fmt.Errorf("You specified wrong encryption handler: OpenSSL for file: %s!", *c.MorphFilename))
	}
	*c.MorphFilename = strings.Split(*c.MorphFilename, ".enc")[0]
	out := fmt.Sprintf("%s/%s", *c.TmpDir, *c.MorphFilename)

	args := []string{
		"aes-256-cbc", "-d", "-a",
		"-md", "sha256",
		"-pbkdf2", "-iter", "20000",
		"-k", *c.Encryptor.OpenSSL.Password,
		"-in", in,
		"-out", out,
	}
	h.RunCommand("openssl", args)

	err := os.Remove(in)
	h.LogErrorExit(err)

	return nil
}

func (o *OpenSSL) Post(c *base.Config) error {
	return nil
}
