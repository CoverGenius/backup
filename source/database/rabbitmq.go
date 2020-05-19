package database

import (
	"bitbucket.org/welovetravel/backup/base"
)

type RabbitMQ struct {
}

func (r *RabbitMQ) Verify(c *base.Config) error {
	return nil
}

func (r *RabbitMQ) Pre(c *base.Config) error {
	return nil
}

func (r *RabbitMQ) Backup(c *base.Config) error {
	return nil
}

func (r *RabbitMQ) Restore(c *base.Config) error {
	return nil
}

func (r *RabbitMQ) Post(c *base.Config) error {
	return nil
}
