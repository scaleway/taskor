package serializer

import (
	"encoding/json"
	"errors"
)

// serializerJson serilize to JSON type
type serializerJSON struct{}

// Serialize value
func (s *serializerJSON) Serialize(data interface{}) ([]byte, error) {
	marshaledData, err := json.Marshal(data)
	if err != nil {
		return nil, errors.New("Failed to encode json")
	}
	return marshaledData, nil
}

// Unserialize value
func (s *serializerJSON) Unserialize(v interface{}, data []byte) error {
	err := json.Unmarshal(data, &v)
	if err != nil {
		return errors.New("Failed to encode json")
	}
	return nil
}
