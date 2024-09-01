package h3test

import (
	"encoding/json"
	"encoding/xml"
)

func ToJSON[T any](v T) ([]byte, error) {
	return json.Marshal(v)
}

func ToXML[T any](v T) ([]byte, error) {
	return xml.Marshal(v)
}

type Stringer interface {
	~string | ~[]byte
}

func ToPlain[T Stringer](v T) []byte {
	return []byte(v)
}
