package base

import (
	"fmt"
	"time"
)

const (
	time_format = "2006-01-02T15-04-05"
)

type AWSMeta struct {
	Account         *string `yaml:"account,omitempty"`
	AccessKey       *string `yaml:"access_key,omitempty"`
	SecretAccessKey *string `yaml:"secret_access_key,omitempty"`
	KMSKeyId        *string `yaml:"kms_key_id,omitempty"`
}

type RDS struct {
	Region               *string `yaml:"region,omitempty"`
	Source               AWSMeta `yaml:"source"`
	Destination          AWSMeta `yaml:"destination"`
	DBInstanceIdentifier *string `yaml:"db_instance_identifier"`
	DBSnapshotSuffix     *string `yaml:"db_snapshot_suffix,omitempty"`
	Snapshots            []*string
}

type S3 struct {
	Bucket  *string `yaml:"bucket,omitempty"`
	Region  *string `yaml:"region,omitempty"`
	AWSMeta AWSMeta `yaml:"meta,omitempty"`
}

type Database struct {
	Redis         Redis         `yaml:"redis,omitempty"`
	Elasticsearch Elasticsearch `yaml:"elasticsearch,omitempty"`
	Postgres      Postgres      `yaml:"postgres,omitempty"`
	MySQL         MySQL         `yaml:"mysql,omitempty"`
	RabbitMQ      RabbitMQ      `yaml:"rabbitmq,omitempty"`
}

type RabbitMQ struct {
	Directory []*string `yaml:"directory"`
}

type BasicAuth struct {
	Enabled  *bool   `yaml:"enabled"`
	Username *string `yaml:"username"`
	Password *string `yaml:"password"`
}

type Elasticsearch struct {
	Host *string   `yaml:"host"`
	Port *string   `yaml:"port"`
	Auth BasicAuth `yaml:"auth,omitempty"`
}

type Postgres struct {
	DBMeta DBMeta    `yaml:"db_connection_settings"`
	Schema []*string `yaml:"schema,omitempty"`
}

type DBMeta struct {
	LoginPath  *string   `yaml:"login_path"`
	Database   []*string `yaml:"database,omitempty"`
	Table      []*string `yaml:"table,omitempty"`
	DumpSchema *bool     `yaml:"dump_schema,omitempty"`
	Options    []*string `yaml:"options,omitempty"`
}

type MySQL struct {
	LoginPath *string  `yaml:"login_path"`
	Databases []string `yaml:"databases,omitempty"`
	Tables    []string `yaml:"tables,omitempty"`
	Schema    *bool    `yaml:"schema,omitempty"`
	Options   []string `yaml:"options,omitempty"`
	Handler   *string  `yaml:"handler,omitempty"`
}

type Redis struct {
	Host *string `yaml:"host,omitempty"`
	Port *string `yaml:"port,omitempty"`
	Auth *string `yaml:"auth,omitempty"`
}

type Directory struct {
	Path []*string `yaml:"path"`
}

type Backup struct {
	Type      *string   `yaml:"type,omitempty"`
	RDS       RDS       `yaml:"rds,omitempty"`
	Database  Database  `yaml:"database,omitempty"`
	S3        S3        `yaml:"s3,omitempty"`
	Directory Directory `yaml:"directory,omitempty"`
}

type TAR struct {
	Options []string `yaml:"options,omitempty"`
}

type Archive struct {
	Type *string `yaml:"type,omitempty"`
	TAR  TAR     `yaml:"tar,omitempty"`
}

type GZIP struct {
	Options []string `yaml:"options,omitempty"`
}

type BZIP2 struct {
	Options []string `yaml:"bzip2,omitempty"`
}

type PIGZ struct {
	Options []string `yaml:"pigz,omitempty"`
}

type Compress struct {
	Type  *string `yaml:"type,omitempty"`
	GZIP  GZIP    `yaml:"gzip,omitempty"`
	BZIP2 BZIP2   `yaml:"bzip2,omitempty"`
	PIGZ  PIGZ    `yaml:"pigz,omitempty"`
}

type Log struct {
	Type *string `yaml:"type,omitempty"`
	Path *string `yaml:"path,omitempty"`
}

type Slack struct {
	WebhookURL *string `yaml:"webhook_url,omitempty"`
	Channel    *string `yaml:"channel,omitempty"`
	Username   *string `yaml:"username,omitempty"`
}

type Mail struct {
	From    *string `yaml:"from,omitempty"`
	To      *string `yaml:"to,omitempty"`
	Address *string `yaml:"address,omitempty"` // host:port
}

type Notification struct {
	Type   []*string `yaml:"type,omitempty"`
	Status Status
	Slack  Slack `yaml:"slack,omitempty"`
	Mail   Mail  `yaml:"mail,omitempty"`
}

type Status struct {
	Status    *string
	StartTime *time.Time
	EndTime   *time.Time
	Duration  *float64
}

func (s *Status) FormatStartTime() *string {
	f := (*s.StartTime).Format(time_format)
	return &f
}

func (s *Status) FormatEndTime() *string {
	f := (*s.EndTime).Format(time_format)
	return &f
}

func (s *Status) CalculateDuration() {
	d := (*s.EndTime).Sub(*s.StartTime).Seconds()
	s.Duration = &d
}

func (s *Status) DurationToString() *string {
	d := fmt.Sprintf("%v seconds", *s.Duration)
	return &d
}

type Split struct {
	Type *string `yaml:"type,omitempty"`
	Size *uint8  `yaml:"size,omitempty"` // in megabytes
}

type OpenSSL struct {
	Password *string `yaml:"password"`
}

type GPG struct {
	Key *string `yaml:"key"`
}

type Encrypt struct {
	Type     *string `yaml:"type,omitempty"`
	GPG      GPG     `yaml:"gpg,omitempty"`
	OpenSSL  OpenSSL `yaml:"openssl,omitempty"`
	Password *string `yaml:"password,omitempty"`
}

type SFTP struct {
	Host      *string `yaml:host",omitempty"`
	Port      *uint8  `yaml:"port,omitempty"`
	Key       *string `yaml:"key,omitempty"`
	Username  *string `yaml:"username,omitempty"`
	Directory *string `yaml:"directory,omitempty"`
}

type Store struct {
	Type *string `yaml:"type,omitempty"`
	Key  *string `yaml:"key,omitempty"`
	SFTP SFTP    `yaml:"sftp,omitempty"`
	S3   S3      `yaml:"s3,omitempty"`
	Path *string `yaml:"path,omitempty"`
}

type Config struct {
	Name             *string      `yaml:"name"`
	Restore          *bool        `yaml:"restore,omitempty"`
	Backup           Backup       `yaml:"backup,omitempty"`
	Archiver         Archive      `yaml:"archiver,omitempty"`
	Compressor       Compress     `yaml:"compressor,omitempty"`
	Notifier         Notification `yaml:"notifier,omitempty"`
	Splitter         Split        `yaml:"splitter,omitempty"`
	Encryptor        Encrypt      `yaml:"encryptor,omitempty"`
	Storage          Store        `yaml:"storage,omitempty"`
	LogLevel         *string      `yaml:"log_level,omitempty"`
	LogFile          *string
	TmpDir           *string `yaml:"tmp_dir,omitempty"`
	TmpDirCleanup    *bool   `yaml:"tmp_dir_cleanup,omitempty"`
	ArchiverSuffix   *string
	CompressorSuffix *string
	EncryptorSuffix  *string
	MorphFilename    *string
	Keep             *uint8 `yaml:"keep,omitempty"`
}
