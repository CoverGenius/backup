package database

import (
	"github.com/CoverGenius/backup/base"
)

type Elasticsearch struct {
}

func (e *Elasticsearch) Verify(c *base.Config) error {
	return nil
}

func (e *Elasticsearch) Pre(c *base.Config) error {
	return nil
}

func (e *Elasticsearch) Backup(c *base.Config) error {
	return nil
}

func (e *Elasticsearch) Restore(c *base.Config) error {
	return nil
}

func (e *Elasticsearch) Post(c *base.Config) error {
	return nil
}
