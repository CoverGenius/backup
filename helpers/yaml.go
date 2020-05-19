package helpers

import (
	"gopkg.in/yaml.v2"
	l "log"
)

func YAMLDecode(f *string, i interface{}) {
	if IsFileExists(f) {
		content := ReadFile(f)
		err := yaml.Unmarshal(*content, i)
		LogErrorExit(err)
	} else {
		l.Fatalf("File: %s does not exists!\n", *f)
	}
}
