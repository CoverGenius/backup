package compressor

import (
	"github.com/CoverGenius/backup/base"
)

type BZIP2 struct {
}

func (b *BZIP2) Verify(c *base.Config) error {
	return nil
}

func (b *BZIP2) Pre(c *base.Config) error {
	return nil
}

func (b *BZIP2) Compress(c *base.Config) error {
	return nil
}

func (b *BZIP2) Decompress(c *base.Config) error {
	return nil
}

func (b *BZIP2) Post(c *base.Config) error {
	return nil
}
