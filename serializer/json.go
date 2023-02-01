package serializer

import (
	"encoding/json"
	"fmt"
)

// serializerJson serilize to JSON type
type serializerJSON struct{}

// Serialize value
func (s *serializerJSON) Serialize(data interface{}) ([]byte, error) {
	marshaledData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to encode json: %v", err)
	}
	return marshaledData, nil
}

// Unserialize value
func (s *serializerJSON) Unserialize(v interface{}, data []byte) error {
	err := json.Unmarshal(data, &v)
	if err != nil {
		return fmt.Errorf("failed to decode json: %v", err)
	}
	return nil
}
