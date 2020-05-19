package helpers

import (
	"os"
	"testing"
)

func TestRunCommand(t *testing.T) {
	command := "touch"
	args := []string{"/tmp/tstcmd"}

	RunCommand(command, args)
	defer os.Remove(args[0])
	finfo, err := os.Stat(args[0])
	LogErrorExit(err)

	if finfo.Name() != "tstcmd" {
		t.Errorf("Failed to run execute command %s.", command)

	}
}
