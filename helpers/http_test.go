package helpers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMakeHTTPRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	}))
	defer func() { ts.Close() }()
	response := MakeHTTPRequest(&ts.URL, "GET", nil, nil, nil, nil, false)
	if response.StatusCode != 200 {
		t.Errorf("Expected status code:200, got %d", response.StatusCode)
	}
}
