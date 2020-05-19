package database

import (
	"github.com/CoverGenius/backup/base"
	h "github.com/CoverGenius/backup/helpers"
	"errors"
	"fmt"
	"os"
)

type MySQL struct {
	LoginPath  string
	SchemaFile string
	DataFile   string
}

func (m *MySQL) Verify(c *base.Config) error {
	if *c.Backup.Database.MySQL.Handler == "mysqldump" && h.IsCommandExists(h.StringP("mysqldump")) == false {
		h.LogErrorExit(errors.New("Please install mysqldump!"))
	}
	if *c.Backup.Database.MySQL.Handler == "innobackupex" && h.IsCommandExists(h.StringP("innobackupex")) == false {
		h.LogErrorExit(errors.New("Please install innobackupex!"))
	}
	// During restore we don't care about databases or tables list
	// because all required information should contain in dump files
	if c.Restore != nil && *c.Restore == true {
		if len(c.Backup.Database.MySQL.Databases) == 0 && len(c.Backup.Database.MySQL.Tables) > 0 {
			h.LogErrorExit(errors.New("You MUST specify database if you have non empty tables list"))
		}
		if len(c.Backup.Database.MySQL.Databases) > 1 && len(c.Backup.Database.MySQL.Tables) > 0 {
			h.LogErrorExit(errors.New("If you have non empty tables list you MUST have only 1 database specified!"))
		}
	}

	return nil
}

func (m *MySQL) Pre(c *base.Config) error {
	m.LoginPath = fmt.Sprintf("--login-path=%s", *c.Backup.Database.MySQL.LoginPath)
	m.SchemaFile = fmt.Sprintf("%s/schema.sql", *c.TmpDir)
	m.DataFile = fmt.Sprintf("%s/data.sql", *c.TmpDir)

	return nil
}

func (m *MySQL) Backup(c *base.Config) error {
	h.Log.Debug("Backup MySQL data ...")

	if c.Backup.Database.MySQL.Schema != nil && *c.Backup.Database.MySQL.Schema == true {
		// In case when ONLY 1 database specified, `schema` flag makes sense
		if len(c.Backup.Database.MySQL.Databases) == 1 {
			args := []string{m.LoginPath}
			args = append(args, c.Backup.Database.MySQL.Options...)
			args = append(args, "--no-data", "-B", c.Backup.Database.MySQL.Databases[0])
			args = append(args, "-r", m.SchemaFile)
			h.RunCommand("mysqldump", args)
		}
	}

	args := []string{}
	args = append(args, m.LoginPath)
	args = append(args, c.Backup.Database.MySQL.Options...)

	// If more than 1 database specified we dump them explicitly
	if len(c.Backup.Database.MySQL.Databases) > 1 {
		args = append(args, "-B")
		args = append(args, c.Backup.Database.MySQL.Databases...)
		// If 0 database specified we dump all databases
	} else if len(c.Backup.Database.MySQL.Databases) == 0 {
		args = append(args, "--all-databases")
		// If only 1 database specified without tables we dump it explicitly
	} else if len(c.Backup.Database.MySQL.Databases) == 1 && len(c.Backup.Database.MySQL.Tables) == 0 {
		args = append(args, "-B", c.Backup.Database.MySQL.Databases[0])
	} else {
		// In any other case we just dump tables for specified database
		args = append(args, c.Backup.Database.MySQL.Databases[0])
		args = append(args, c.Backup.Database.MySQL.Tables...)
	}
	args = append(args, "-r", m.DataFile)
	h.RunCommand("mysqldump", args)

	// If only 1 database specified and tables list not empty, mysqldump
	// does not add `USE "${database}";` instruction into dump file.
	// In order to simplify restore logic next block will add this instruction explicitly
	if len(c.Backup.Database.MySQL.Databases) == 1 && len(c.Backup.Database.MySQL.Tables) > 0 {
		use_instruction := []byte(fmt.Sprintf("USE `%s`;\n-- ", c.Backup.Database.MySQL.Databases[0]))
		f, err := os.OpenFile(m.DataFile, os.O_RDWR, 0600)
		defer f.Close()
		h.LogErrorExit(err)

		_, err = f.WriteAt(use_instruction, 0)
		h.LogErrorExit(err)
	}

	return nil
}

func (m *MySQL) Restore(c *base.Config) error {
	h.Log.Debug("Restore MySQL data ...")

	restoreSchemaCmd := fmt.Sprintf("source %s;", m.SchemaFile)
	restoreDataCmd := fmt.Sprintf("source %s;", m.DataFile)

	if *c.Backup.Database.MySQL.Schema == true {
		args := []string{
			m.LoginPath, "information_schema",
			"-e", restoreSchemaCmd,
		}
		h.RunCommand("mysql", args)
	}

	args := []string{
		m.LoginPath, "information_schema",
		"-e", restoreDataCmd,
	}
	h.RunCommand("mysql", args)

	return nil
}

func (m *MySQL) Post(c *base.Config) error {
	return nil
}
