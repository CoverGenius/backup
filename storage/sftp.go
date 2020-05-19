package storage

import (
	"github.com/CoverGenius/backup/base"
)

type SFTP struct {
}

func (s *SFTP) Verify(c *base.Config) error {
	return nil
}

func (s *SFTP) Pre(c *base.Config) error {
	return nil
}

func (s *SFTP) Store(c *base.Config) error {
	return nil
}

func (s *SFTP) Fetch(c *base.Config) error {
	return nil
}

func (s *SFTP) Post(c *base.Config) error {
	return nil
}
