package helpers

import (
	"bytes"
	"encoding/json"
	"io"
)

func ConvertBodyToString(body io.ReadCloser) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(body) //nolint
	return buf.String()
}

func MarshalResponse(data interface{}) string {
	res, _ := json.Marshal(data)
	return string(res)
}
