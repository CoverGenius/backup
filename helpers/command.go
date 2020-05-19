package helpers

import (
	"os"
	"os/exec"
)

func RunCommand(name string, arg []string) {
	path, err := exec.LookPath(name)
	LogErrorExit(err)

	cmd := exec.Command(path, arg...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	LogErrorExit(err)
}
