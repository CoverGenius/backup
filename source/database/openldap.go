package database

import (
	"github.com/CoverGenius/backup/base"
)

type OpenLDAP struct {
}

func (o *OpenLDAP) Verify(c *base.Config) error {
	return nil
}

func (o *OpenLDAP) Pre(c *base.Config) error {
	return nil
}

func (o *OpenLDAP) Backup(c *base.Config) error {
	return nil
}

func (o *OpenLDAP) Restore(c *base.Config) error {
	return nil
}

func (o *OpenLDAP) Post(c *base.Config) error {
	return nil
}
