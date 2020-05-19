package helpers

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

func ReadFile(path *string) *[]byte {
	absolutePath, err := filepath.Abs(*path)
	LogError(err)
	content, err := ioutil.ReadFile(absolutePath)
	LogError(err)
	return &content
}

func IsFileExists(path *string) bool {
	_, err := os.Stat(*path)
	if err == nil {
		return true
	} else {
		return false
	}
}

func IsCommandExists(command *string) bool {
	_, err := exec.LookPath(*command)
	if err == nil {
		return true
	} else {
		return false
	}
}

func RemoveDir(path *string) {
	absolutePath, err := filepath.Abs(*path)
	LogError(err)

	matched, err := regexp.Match(
		`^(/|/etc|/var|/sys|/dev|/boot|/proc|/lib|/lib64|/root|/run|/tmp|/usr)$`,
		[]byte(absolutePath),
	)

	LogErrorExit(err)
	if matched == true {
		LogErrorExit(fmt.Errorf("Cannot delete system directory: %s!", absolutePath))
	}

	os.RemoveAll(absolutePath)
}
