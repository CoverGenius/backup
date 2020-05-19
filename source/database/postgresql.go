package database

import (
	"github.com/CoverGenius/backup/base"
)

type Postgres struct {
}

func (p *Postgres) Verify(c *base.Config) error {
	return nil
}

func (p *Postgres) Pre(c *base.Config) error {
	return nil
}

func (p *Postgres) Backup(c *base.Config) error {
	return nil
}

func (p *Postgres) Restore(c *base.Config) error {
	return nil
}

func (p *Postgres) Post(c *base.Config) error {
	return nil
}
