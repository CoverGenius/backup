package helpers

import (
	"testing"
)

func TestFormatString(t *testing.T) {
	format := "{KEY1}/{KEY2}"
	expected := "abc/123"
	result := FormatString(&format, "{KEY1}", "abc", "{KEY2}", "123")
	if result != expected {
		t.Errorf("Result is incorrect, got: %s, want: %s.", result, expected)
	}
}
