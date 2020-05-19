package splitter

import (
	"github.com/CoverGenius/backup/base"
)

type Split struct {
}

func (s *Split) Verify(c *base.Config) error {
	return nil
}

func (s *Split) Pre(c *base.Config) error {
	return nil
}

func (s *Split) Split(c *base.Config) error {
	return nil
}

func (s *Split) Post(c *base.Config) error {
	return nil
}
