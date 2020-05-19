package archiver

import (
	"github.com/CoverGenius/backup/base"
	h "github.com/CoverGenius/backup/helpers"
	"errors"
	"fmt"
	"os"
	"strings"
)

type TAR struct {
}

func (t *TAR) Verify(c *base.Config) error {
	if h.IsCommandExists(h.StringP("tar")) == false {
		h.LogErrorExit(errors.New("Please install tar!"))
	}

	return nil
}

func (t *TAR) Pre(c *base.Config) error {
	return nil
}

func (t *TAR) Archive(c *base.Config) error {
	h.Log.Debug("Archive data using Tar ...")

	fPath := fmt.Sprintf("%s/backup.%s", *c.TmpDir, *c.ArchiverSuffix)

	args := []string{"-C", *c.TmpDir}
	args = append(args, c.Archiver.TAR.Options...)
	args = append(args, "--exclude=backup.tar", "-cvf", fPath, ".")

	h.RunCommand("tar", args)

	return nil
}

func (t *TAR) Unarchive(c *base.Config) error {
	h.Log.Debug("Unarchive data using Tar ...")

	parts := strings.Split(*c.MorphFilename, ".")
	if *(h.GetLastStringElement(parts)) != "tar" {
		h.LogErrorExit(fmt.Errorf("You specified wrong archive handler: TAR for file: %s!", *c.MorphFilename))
	}

	fPath := fmt.Sprintf("%s/%s", *c.TmpDir, *c.MorphFilename)
	args := []string{"-C", *c.TmpDir, "-xvf", fPath}

	h.RunCommand("tar", args)

	err := os.Remove(fPath)
	h.LogErrorExit(err)

	return nil
}

func (t *TAR) Post(c *base.Config) error {
	return nil
}
