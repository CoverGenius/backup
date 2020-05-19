package helpers

import (
	"bytes"
	"crypto/tls"
	"net/http"
	h "net/http"
	u "net/url"
)

func MakeHTTPRequest(urlp *string, method string, headers *map[string]string, username *string, password *string, body *bytes.Buffer, skipTlsVerification bool) *h.Response {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipTlsVerification},
	}
	client := &h.Client{Transport: tr}
	url, err := u.Parse(*urlp)
	LogError(err)
	var req *h.Request
	if body == nil {
		req, err = h.NewRequest(method, url.String(), nil)
	} else {
		req, err = h.NewRequest(method, url.String(), body)
	}
	LogError(err)
	if username != nil && password != nil {
		req.SetBasicAuth(*username, *password)
	}
	if headers != nil && len(*headers) > 0 {
		for in, elem := range *headers {
			req.Header.Add(in, elem)
		}
	} else {
		req.Header.Add("Content-Type", "application/json")
	}
	resp, err := client.Do(req)
	LogError(err)
	return resp
}
