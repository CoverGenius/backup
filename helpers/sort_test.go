package helpers

import (
	"testing"
	"time"
)

func TestTDataQuickSort(t *testing.T) {
	tu := time.Now().Unix()
	expected := `this is correct sorting order for TData data structure ... .. . ! `
	testdata := []*TData{
		{Timestamp: tu + 100, Data: "TData"},
		{Timestamp: tu - 300, Data: ".."},
		{Timestamp: tu - 400, Data: "."},
		{Timestamp: tu + 600, Data: "is"},
		{Timestamp: tu + 700, Data: "this"},
		{Timestamp: tu + 150, Data: "correct"},
		{Timestamp: tu - 300, Data: "..."},
		{Timestamp: tu + 50, Data: "data"},
		{Timestamp: tu - 700, Data: "!"},
		{Timestamp: tu + 123, Data: "for"},
		{Timestamp: tu + 123, Data: "order"},
		{Timestamp: tu + 124, Data: "sorting"},
		{Timestamp: tu + 0, Data: "structure"},
	}
	TDataQuickSort(testdata)
	var result string
	for _, z := range testdata {
		result = result + (*z).Data + " "
	}
	if result != expected {
		t.Errorf("Wrong result. Got: %s, expected: %s", result, expected)
	}
}
