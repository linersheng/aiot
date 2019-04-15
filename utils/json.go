package utils

import (
	"github.com/json-iterator/go"
	//"encoding/json"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func JsonUnmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func JsonMarshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
