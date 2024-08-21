package test

import "encoding/json"

func MarshalResponse(data interface{}) string {
	res, _ := json.Marshal(data)
	return string(res)
}
