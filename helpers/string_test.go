package helpers

import (
	"testing"
)

func TestGetLastStringElement(t *testing.T) {
	data := []string{"123", "234", "abc"}
	expected := "abc"
	result := GetLastStringElement(data)

	if *result != expected {
		t.Errorf("Result is incorrect, got: %s, want: %s.", *result, expected)
	}
}
