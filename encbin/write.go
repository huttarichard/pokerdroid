package encbin

import (
	"encoding/binary"
	"io"
	"reflect"
)

type BinaryWriter interface {
	WriteBinary(w io.Writer) error
	Size() uint64
}

// WriteWithLen writes the size header followed by direct WriteBinary call
// This avoids full in-memory marshaling of large nested structures
func WriteWithLen(w io.Writer, v BinaryWriter) error {
	// Handle nil values
	if v == nil {
		return MarshalValues(w, ^uint64(0))
	}

	vx := reflect.ValueOf(v)
	if vx.Kind() == reflect.Ptr || vx.Kind() == reflect.Interface {
		if vx.IsNil() {
			return binary.Write(w, binary.LittleEndian, ^uint64(0))
		}
	}

	// Write size header
	size := v.Size()
	if err := MarshalValues(w, size); err != nil {
		return err
	}

	// Directly write the binary data
	return v.WriteBinary(w)
}
