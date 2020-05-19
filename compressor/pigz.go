package compressor

import (
	"bitbucket.org/welovetravel/backup/base"
)

type PIGZ struct {
}

func (p *PIGZ) Verify(c *base.Config) error {
	return nil
}

func (p *PIGZ) Pre(c *base.Config) error {
	return nil
}

func (p *PIGZ) Compress(c *base.Config) error {
	return nil
}

func (p *PIGZ) Decompress(c *base.Config) error {
	return nil
}

func (p *PIGZ) Post(c *base.Config) error {
	return nil
}
