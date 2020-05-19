package compressor

import (
	"bitbucket.org/welovetravel/backup/base"
	h "bitbucket.org/welovetravel/backup/helpers"
	"errors"
	"fmt"
	"strings"
)

type GZIP struct {
}

func (g *GZIP) Verify(c *base.Config) error {
	if h.IsCommandExists(h.StringP("gzip")) == false {
		h.LogErrorExit(errors.New("Please install gzip!"))
	}

	return nil
}

func (g *GZIP) Pre(c *base.Config) error {
	return nil
}

func (g *GZIP) Compress(c *base.Config) error {
	h.Log.Debug("Compressing data using GZIP ...")

	fPath := fmt.Sprintf("%s/backup.%s", *c.TmpDir, *c.ArchiverSuffix)

	args := []string{}
	args = append(args, c.Compressor.GZIP.Options...)
	args = append(args, fPath)

	h.RunCommand("gzip", args)

	return nil
}

func (g *GZIP) Decompress(c *base.Config) error {
	h.Log.Debug("Uncompressing data using GZIP ...")

	parts := strings.Split(*c.MorphFilename, ".")
	if *(h.GetLastStringElement(parts)) != "gz" {
		h.LogErrorExit(fmt.Errorf("You specified wrong compression handler: GZIP for file: %s!", *c.MorphFilename))
	}

	fPath := fmt.Sprintf("%s/%s", *c.TmpDir, *c.MorphFilename)
	args := []string{
		"-d", fPath,
	}

	h.RunCommand("gzip", args)

	*c.MorphFilename = strings.Split(*c.MorphFilename, ".gz")[0]

	return nil
}

func (g *GZIP) Post(c *base.Config) error {
	return nil
}
