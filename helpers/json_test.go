package helpers

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testJSON struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func TestJSONEncode(t *testing.T) {
	expected := "{\"name\":\"test\",\"value\":123}"
	var jt testJSON
	var b bytes.Buffer
	jt.Name = "test"
	jt.Value = 123
	JSONEncode(&b, &jt)
	result := b.String()
	if result != expected {
		t.Errorf("Expected string: %s, got: %s", expected, result)
	}
}

func TestJSONDecode(t *testing.T) {
	var jt testJSON
	handler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "{\"name\":\"test\",\"value\":123}")
	}
	req := httptest.NewRequest("GET", "http://example.com/test", nil)
	w := httptest.NewRecorder()
	handler(w, req)
	resp := w.Result()
	JSONDecode(resp, &jt)
	if jt.Name != "test" || jt.Value != 123 {
		t.Errorf("Expected name:test and value:123, got name:%s and value:%d", jt.Name, jt.Value)
	}
}
