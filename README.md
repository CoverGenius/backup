# Backup tool

Tool used to backup and restore things.
Currently only supports RDS and MySQL databases.


## Example usage for RDS backup

Run a one-time backup/restore:

```
./backup -config </path/to/your/config.yaml>
```

## Config file examples you can find in `examples` directory.

# Permissions and setup

Minimal AWS IAM permissions required to perform RDS Snapshot backup
```
---
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "RDSBackup",
      "Effect": "Allow",
      "Action": [
        "rds:DescribeDBInstances",
        "rds:CopyDBSnapshot",
        "rds:DescribeDBSnapshots",
        "rds:DescribeDBSnapshotAttributes",
        "rds:DeleteDBSnapshot",
        "rds:CreateDBSnapshot",
        "rds:ModifyDBSnapshotAttribute",
        "kms:DescribeKey",
        "kms:CreateGrant"
      ],
      "Resource": "*"
    }
  ]
}
```
# In order to work with MySQL you MUST configure and use login-path
```
mysql_config_editor set --host=<host> --port=3306 --login-path=<name> --user=<user> --password
<you will be prompted for password>
```
Basically above command will create credentials file which you can refer to in order to connect to MySQL.


# Future plans

Support the following backup sources:

- Elasticsearch
- PostgreSQL
- Redis
- RabbitMQ
..
