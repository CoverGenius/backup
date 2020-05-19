package base

type Source interface {
	Verify(c *Config) error
	Pre(c *Config) error
	Backup(c *Config) error
	Restore(c *Config) error
	Post(c *Config) error
}

type Archiver interface {
	Verify(c *Config) error
	Pre(c *Config) error
	Archive(c *Config) error
	Unarchive(c *Config) error
	Post(c *Config) error
}

type Compressor interface {
	Verify(c *Config) error
	Pre(c *Config) error
	Compress(c *Config) error
	Decompress(c *Config) error
	Post(c *Config) error
}

type Logger interface {
	Verify(c *Config) error
	Pre(c *Config) error
	Configure(c *Config) error
	Post(c *Config) error
}

type Notifier interface {
	Verify(c *Config) error
	Pre(c *Config) error
	Notify(c *Config) error
	Post(c *Config) error
}

type Storage interface {
	Verify(c *Config) error
	Pre(c *Config) error
	Store(c *Config) error
	Fetch(c *Config) error
	Post(c *Config) error
}

type Encryptor interface {
	Verify(c *Config) error
	Pre(c *Config) error
	Encrypt(c *Config) error
	Decrypt(c *Config) error
	Post(c *Config) error
}

type Cleaner interface {
	Verify(c *Config) error
	Pre(c *Config) error
	Cleanup(c *Config) error
	Post(c *Config) error
}

type Splitter interface {
	Verify(c *Config) error
	Pre(c *Config) error
	Split(c *Config) error
	Post(c *Config) error
}
