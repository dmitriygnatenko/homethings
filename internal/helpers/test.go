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

func ConvertDTOToIOReader(data interface{}) io.Reader {
	var res io.Reader
	if data != nil {
		bytesData, _ := json.Marshal(data)
		res = bytes.NewReader(bytesData)
	}

	return res
}

func MarshalResponse(data interface{}) string {
	res, _ := json.Marshal(data)
	return string(res)
}
