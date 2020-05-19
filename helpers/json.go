package helpers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func JSONEncode(b *bytes.Buffer, i interface{}) {
	j, err := json.Marshal(i)
	LogError(err)
	err = json.Compact(b, j)
	LogError(err)
}

func JSONDecode(r *http.Response, i interface{}) {
	result, err := ioutil.ReadAll(r.Body)
	LogError(err)
	json.Unmarshal(result, i)
}
