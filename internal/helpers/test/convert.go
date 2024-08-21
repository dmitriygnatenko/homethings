package test

import (
	"bytes"
	"encoding/json"
	"io"
)

func ConvertBodyToString(body io.ReadCloser) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(body) // nolint
	return buf.String()
}

func ConvertDataToIOReader(data interface{}) io.Reader {
	var res io.Reader
	if data != nil {
		bytesData, _ := json.Marshal(data)
		res = bytes.NewReader(bytesData)
	}

	return res
}
