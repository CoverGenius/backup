package encryptor

import (
	"bitbucket.org/welovetravel/backup/base"
)

type GPG struct {
}

func (g *GPG) Verify(c *base.Config) error {
	return nil
}

func (g *GPG) Pre(c *base.Config) error {
	return nil
}

func (g *GPG) Encrypt(c *base.Config) error {
	return nil
}

func (g *GPG) Decrypt(c *base.Config) error {
	return nil
}

func (g *GPG) Post(c *base.Config) error {
	return nil
}
