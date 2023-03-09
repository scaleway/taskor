package serializer

// Type type
type Type int

// Serializer constant
const (
	TypeJSON Type = 1 + iota
	TypeGob
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
	return GetSerializer(GlobalSerializer)
}

// GetSerializer Return serilizer from type
func GetSerializer(Type Type) Serializer {
	var serial Serializer
	switch {
	case Type == TypeJSON:
		serial = &serializerJSON{}
	case Type == TypeGob:
		serial = &serializerGob{}
	default:
		serial = &serializerJSON{}
	}
	return serial
}

func GetContentType(Type Type) string {
	switch {
	case Type == TypeJSON:
		return "text/plain"
	case Type == TypeGob:
		return "application/octet-stream"
	default:
		return "text/plain"
	}
}
