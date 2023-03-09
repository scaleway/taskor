package serializer

import (
	"encoding/gob"
	"testing"

	"github.com/stretchr/testify/require"
)

type TestInterface interface {
	dummy()
}

type TestInterfaceImpl1 struct {
	S string
}

func (t TestInterfaceImpl1) dummy() {}

type TestInterfaceImpl2 struct {
	I int
}

func (t TestInterfaceImpl2) dummy() {}

type TestStruct struct {
	I int
	S string
	A []string
	T TestInterface
	P *string
}

func Test_serializerGob(t *testing.T) {
	for _, concreteType := range []any{TestInterfaceImpl1{}, TestInterfaceImpl2{}} {
		gob.Register(concreteType)
	}

	testString := "hello"

	tests := []struct {
		name string
		data any
	}{
		{
			name: "without any attributes set",
			data: TestStruct{},
		},
		{
			name: "with some attributes with scalar type set",
			data: TestStruct{
				I: 1,
			},
		},
		{
			name: "with some attributes with scalar type set",
			data: TestStruct{
				S: "hello",
			},
		},
		{
			name: "with some attributes with array type set",
			data: TestStruct{
				A: []string{"hello", "world"},
			},
		},
		{
			name: "with some attributes with interface type set",
			data: TestStruct{
				T: TestInterfaceImpl1{
					S: "hello",
				},
			},
		},
		{
			name: "with some attributes with interface type set",
			data: TestStruct{
				T: TestInterfaceImpl2{
					I: 1,
				},
			},
		},
		{
			name: "with some attributes with pointer type set",
			data: TestStruct{
				P: &testString,
			},
		},
		{
			name: "with all attributes set",
			data: TestStruct{
				I: 1,
				S: "hello",
				A: []string{"hello", "world"},
				T: TestInterfaceImpl1{
					S: "hello",
				},
				P: &testString,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serializer := &serializerGob{}

			serialized, err := serializer.Serialize(tt.data)
			require.NoError(t, err)

			var unserialized TestStruct
			err = serializer.Unserialize(&unserialized, serialized)
			require.NoError(t, err)

			require.Equal(t, tt.data, unserialized)
		})
	}
}
