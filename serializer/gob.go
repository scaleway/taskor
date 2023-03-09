package serializer

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type serializerGob struct{}

func (s *serializerGob) Serialize(data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(data)
	if err != nil {
		return nil, fmt.Errorf("failed to encode with gob: %v", err)
	}
	return buf.Bytes(), nil
}

func (s *serializerGob) Unserialize(v interface{}, data []byte) error {
	buf := bytes.NewReader(data)
	err := gob.NewDecoder(buf).Decode(v)
	if err != nil {
		return fmt.Errorf("failed to decode with gob: %v", err)
	}
	return nil
}
