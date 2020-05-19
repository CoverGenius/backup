package storage

import (
	"github.com/CoverGenius/backup/base"
)

type Local struct {
}

func (l *Local) Verify(c *base.Config) error {
	return nil
}

func (l *Local) Pre(c *base.Config) error {
	return nil
}

func (l *Local) Store(c *base.Config) error {
	return nil
}

func (l *Local) Fetch(c *base.Config) error {
	return nil
}

func (l *Local) Post(c *base.Config) error {
	return nil
}
