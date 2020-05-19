package database

import (
	"bitbucket.org/welovetravel/backup/base"
	h "bitbucket.org/welovetravel/backup/helpers"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/rds"
	"sync"
)

type Session struct {
	Source      *rds.RDS
	Destination *rds.RDS
}

type RDS struct {
	Session Session
}

func (r *RDS) Verify(c *base.Config) error {
	if c.Backup.RDS.Source.AccessKey == nil || c.Backup.RDS.Source.SecretAccessKey == nil {
		return errors.New("Access key or Secret Key is wrong for source!")
	}
	if c.Backup.RDS.Destination.AccessKey == nil || c.Backup.RDS.Destination.SecretAccessKey == nil {
		return errors.New("Access key or Secret Key is wrong for destination!")
	}
	if c.Backup.RDS.DBInstanceIdentifier == nil {
		return errors.New("No DB instance identifier specified")
	}
	if c.Backup.RDS.Destination.KMSKeyId == nil {
		return errors.New("KMS key MUST be specified for destination")
	}
	if c.Backup.RDS.Destination.Account == nil || c.Backup.RDS.Source.Account == nil {
		return errors.New("Account ID MUST be specified for both source and destination")
	}
	return nil
}

func (r *RDS) Pre(c *base.Config) error {
	r.Session.Source = rds.New(session.Must(session.NewSession(&aws.Config{
		Region:      c.Backup.RDS.Region,
		Credentials: credentials.NewStaticCredentials(*c.Backup.RDS.Source.AccessKey, *c.Backup.RDS.Source.SecretAccessKey, ""),
	})))
	r.Session.Destination = rds.New(session.Must(session.NewSession(&aws.Config{
		Region:      c.Backup.RDS.Region,
		Credentials: credentials.NewStaticCredentials(*c.Backup.RDS.Destination.AccessKey, *c.Backup.RDS.Destination.SecretAccessKey, ""),
	})))
	return nil
}

func CreateDBSnapshot(dbi *string, si *string, r *rds.RDS) *rds.DBSnapshot {
	source_db_snapshot_input := &rds.CreateDBSnapshotInput{
		DBInstanceIdentifier: dbi,
		DBSnapshotIdentifier: si,
	}
	source_db_snapshot_output, err := r.CreateDBSnapshot(source_db_snapshot_input)
	h.LogErrorExit(err)
	return source_db_snapshot_output.DBSnapshot
}

func CopyDBSnapshot(source *string, target *string, key *string, tags *bool, r *rds.RDS) *rds.DBSnapshot {
	copy_db_snapshot_input := &rds.CopyDBSnapshotInput{
		CopyTags:                   tags,
		TargetDBSnapshotIdentifier: target,
		SourceDBSnapshotIdentifier: source,
		KmsKeyId:                   key,
	}
	copy_db_snapshot_output, err := r.CopyDBSnapshot(copy_db_snapshot_input)
	h.LogErrorExit(err)
	return copy_db_snapshot_output.DBSnapshot
}

func ShareDBSnapshot(s *string, a *string, r *rds.RDS) *rds.DBSnapshotAttributesResult {
	modify_snapshot_input := &rds.ModifyDBSnapshotAttributeInput{
		DBSnapshotIdentifier: s,
		AttributeName:        aws.String("restore"),
		ValuesToAdd: []*string{
			a,
		},
	}
	modify_snapshot_output, err := r.ModifyDBSnapshotAttribute(modify_snapshot_input)
	h.LogErrorExit(err)
	return modify_snapshot_output.DBSnapshotAttributesResult
}

func DeleteDBSnapshot(s *string, wg *sync.WaitGroup, r *rds.RDS) {
	describe_snapshot_input := &rds.DescribeDBSnapshotsInput{
		DBSnapshotIdentifier: s,
	}
	delete_snapshot_input := &rds.DeleteDBSnapshotInput{
		DBSnapshotIdentifier: s,
	}
	_, err := r.DeleteDBSnapshot(delete_snapshot_input)
	h.LogError(err)
	err = r.WaitUntilDBSnapshotDeleted(describe_snapshot_input)
	h.LogError(err)

	wg.Done()
}

func DescribeKey(snapshot *rds.DBSnapshot, c *base.Config) *kms.DescribeKeyOutput {
	kms_session := kms.New(session.Must(session.NewSession(&aws.Config{
		Region:      c.Backup.RDS.Region,
		Credentials: credentials.NewStaticCredentials(*c.Backup.RDS.Source.AccessKey, *c.Backup.RDS.Source.SecretAccessKey, ""),
	})))
	source_kms_input := kms.DescribeKeyInput{
		KeyId: snapshot.KmsKeyId,
	}
	source_kms_output, err := kms_session.DescribeKey(&source_kms_input)
	h.LogErrorExit(err)
	return source_kms_output
}

func SetSnapshotId(c *base.Config) error {
	var snapshotId string
	if c.Backup.RDS.DBSnapshotSuffix != nil {
		snapshotId = fmt.Sprintf("%s-%s-%s", *c.Backup.RDS.DBInstanceIdentifier, *c.Backup.RDS.DBSnapshotSuffix, *h.TimeFormat(h.GetTimeNow()))
	} else {
		snapshotId = fmt.Sprintf("%s-%s", *c.Backup.RDS.DBInstanceIdentifier, *h.TimeFormat(h.GetTimeNow()))
	}
	snapshotIdTmp := fmt.Sprintf("%s-tmp", snapshotId)
	c.Backup.RDS.Snapshots = []*string{
		&snapshotId,
		&snapshotIdTmp,
	}
	return nil
}

func (r *RDS) Backup(c *base.Config) error {
	source_db_instance_input := &rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: c.Backup.RDS.DBInstanceIdentifier,
	}

	source_db_instance_output, err := r.Session.Source.DescribeDBInstances(source_db_instance_input)
	h.LogErrorExit(err)

	if len(source_db_instance_output.DBInstances) == 0 {
		h.LogErrorExit(errors.New("No DB instances found!"))
	}

	SetSnapshotId(c)
	sourceSnapshotPrefix := fmt.Sprintf("arn:aws:rds:%s:%s:snapshot", *c.Backup.RDS.Region, *c.Backup.RDS.Source.Account)

	h.Log.Info("Creating snapshot")
	source_db_snapshot_output := CreateDBSnapshot(c.Backup.RDS.DBInstanceIdentifier, c.Backup.RDS.Snapshots[0], r.Session.Source)
	describe_source_db_snapshot_input := &rds.DescribeDBSnapshotsInput{
		DBSnapshotIdentifier: source_db_snapshot_output.DBSnapshotIdentifier,
	}
	snapshot_to_share := c.Backup.RDS.Snapshots[0]

	h.Log.Info("Wait until source snapshot will be ready!")
	err = r.Session.Source.WaitUntilDBSnapshotAvailable(describe_source_db_snapshot_input)
	h.LogErrorExit(err)

	h.Log.Debug("Checking if snapshot is encrypted ...")
	if *source_db_instance_output.DBInstances[0].StorageEncrypted == true {
		source_kms_output := DescribeKey(source_db_snapshot_output, c)

		if *source_kms_output.KeyMetadata.KeyManager == "AWS" {
			h.Log.Warn("Source snapshot uses default AWS key!")
			h.Log.Warn("Creating another snapshot encrypted with non-default KMS key!")
			copy_db_snapshot_output := CopyDBSnapshot(c.Backup.RDS.Snapshots[0], c.Backup.RDS.Snapshots[1], c.Backup.RDS.Source.KMSKeyId, aws.Bool(true), r.Session.Source)
			h.Log.Info("Wait until snapshot will be available")
			db_snapshot_input := &rds.DescribeDBSnapshotsInput{
				DBSnapshotIdentifier: copy_db_snapshot_output.DBSnapshotIdentifier,
			}
			err = r.Session.Source.WaitUntilDBSnapshotAvailable(db_snapshot_input)
			h.LogErrorExit(err)
			snapshot_to_share = c.Backup.RDS.Snapshots[1]
		}
	}

	h.Log.Info("Sharing snapshot with another account!")
	ShareDBSnapshot(snapshot_to_share, c.Backup.RDS.Destination.Account, r.Session.Source)

	h.Log.Info("Copying shared snapshot to be indepenent from source account")
	sourceSnapshotPrefix = fmt.Sprintf("%s:%s", sourceSnapshotPrefix, *snapshot_to_share)
	copy_db_snapshot_output := CopyDBSnapshot(&sourceSnapshotPrefix, c.Backup.RDS.Snapshots[0], c.Backup.RDS.Destination.KMSKeyId, aws.Bool(false), r.Session.Destination)
	db_snapshot_input := &rds.DescribeDBSnapshotsInput{
		DBSnapshotIdentifier: copy_db_snapshot_output.DBSnapshotIdentifier,
	}
	h.Log.Info("Waiting until shared snapshot will be copied ...")
	err = r.Session.Destination.WaitUntilDBSnapshotAvailable(db_snapshot_input)
	h.LogErrorExit(err)

	return nil
}

func (r *RDS) Restore(c *base.Config) error {
	return nil
}

func (r *RDS) Post(c *base.Config) error {
	var wg sync.WaitGroup

	h.Log.Info("Delete tmp snapshots on source account!")
	for _, snapshot := range c.Backup.RDS.Snapshots {
		wg.Add(1)
		go DeleteDBSnapshot(snapshot, &wg, r.Session.Source)
	}

	h.Log.Info("Delete snapshots on destination account(applying keep policy)")
	describe_db_snapshot_input := &rds.DescribeDBSnapshotsInput{
		DBInstanceIdentifier: c.Backup.RDS.DBInstanceIdentifier,
	}
	dscribe_db_snapshot_output, err := r.Session.Destination.DescribeDBSnapshots(describe_db_snapshot_input)
	h.LogError(err)

	toDelete := []*h.TData{}
	if len(dscribe_db_snapshot_output.DBSnapshots) > int(*c.Keep) {
		// quicksort
		for _, snapshot := range dscribe_db_snapshot_output.DBSnapshots {
			t := &h.TData{
				Timestamp: (*snapshot.SnapshotCreateTime).Unix(),
				Data:      *snapshot.DBSnapshotIdentifier,
			}
			toDelete = append(toDelete, t)
		}
	}
	h.TDataQuickSort(toDelete)
	if len(toDelete) > int(*c.Keep) {
		for _, snapshot := range toDelete[*c.Keep:] {
			h.Log.Info("Delete snapshot: ", snapshot.Data)
			wg.Add(1)
			go DeleteDBSnapshot(&snapshot.Data, &wg, r.Session.Destination)
		}
	}
	wg.Wait()

	return nil
}
