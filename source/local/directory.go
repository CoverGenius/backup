package local

import (
	"github.com/CoverGenius/backup/base"
)

type Directory struct {
	Path []*string
}

func (d *Directory) Verify(c *base.Config) error {
	return nil
}

func (d *Directory) Pre(c *base.Config) error {
	return nil
}

func (d *Directory) Backup(c *base.Config) error {
	return nil
}

func (d *Directory) Restore(c *base.Config) error {
	return nil
}

func (d *Directory) Post(c *base.Config) error {
	return nil
}
