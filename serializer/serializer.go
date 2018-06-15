package serializer

// Type type
type Type int

// Serializer constant
const (
	TypeJSON Type = 1 + iota
)

// GlobalSerializer var use to choose Serializer, should be init
var GlobalSerializer Type

// Serializer engine interface
type Serializer interface {
	Serialize(data interface{}) ([]byte, error)
	Unserialize(v interface{}, data []byte) error
}

// GetGlobalSerializer Return serilizer from type
func GetGlobalSerializer() Serializer {
	var serial Serializer
	switch {
	case GlobalSerializer == TypeJSON:
		serial = &serializerJSON{}
	default:
		serial = &serializerJSON{}
	}
	return serial
}

// GetSerializer Return serilizer from type
func GetSerializer(Type Type) Serializer {
	var serial Serializer
	switch {
	case Type == TypeJSON:
		serial = &serializerJSON{}
	default:
		serial = &serializerJSON{}
	}
	return serial
}
