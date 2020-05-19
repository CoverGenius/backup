package database

import (
	"bitbucket.org/welovetravel/backup/base"
)

type Redis struct {
}

func (t *Redis) Verify(c *base.Config) error {
	return nil
}

func (t *Redis) Pre(c *base.Config) error {
	return nil
}

func (t *Redis) Backup(c *base.Config) error {
	return nil
}

func (t *Redis) Restore(c *base.Config) error {
	return nil
}

func (t *Redis) Post(c *base.Config) error {
	return nil
}
