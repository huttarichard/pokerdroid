package encbin

import (
	"encoding"
	"encoding/binary"
	"errors"
	"io"
	"reflect"
)

type Sizer interface {
	Size() uint64
}

type Size interface {
	uint8 | uint16 | uint32 | uint64
}

func MarshalValues(buf io.Writer, rr ...any) error {
	for _, rr := range rr {
		err := binary.Write(buf, binary.LittleEndian, rr)
		if err != nil {
			return err
		}
	}
	return nil
}

func UnmarshalValues(buf io.Reader, rr ...any) error {
	for _, rr := range rr {
		err := binary.Read(buf, binary.LittleEndian, rr)
		if err != nil {
			return err
		}
	}
	return nil
}

// Add helper to get max value for Sizer type
func maxValue[T Size]() T {
	return ^T(0)
}

func MarshalMap[K comparable, V comparable, S Size](buf io.Writer, rr map[K]V) error {
	// Handle nil map
	if rr == nil {
		return binary.Write(buf, binary.LittleEndian, maxValue[S]())
	}

	length := S(len(rr))
	// Check for length conflict with nil marker
	if length == maxValue[S]() {
		return errors.New("map length equals nil marker")
	}

	err := binary.Write(buf, binary.LittleEndian, length)
	if err != nil {
		return err
	}

	// Write the map entries
	for k, v := range rr {
		err := binary.Write(buf, binary.LittleEndian, k)
		if err != nil {
			return err
		}
		err = binary.Write(buf, binary.LittleEndian, v)
		if err != nil {
			return err
		}
	}

	return nil
}

func UnmarshalMap[K comparable, V comparable, S Size](buf io.Reader) (map[K]V, error) {
	var length S

	err := binary.Read(buf, binary.LittleEndian, &length)
	if err != nil {
		return nil, err
	}

	// Check for nil marker
	if length == maxValue[S]() {
		return nil, nil
	}

	rr := make(map[K]V, length)

	for i := S(0); i < length; i++ {
		var k K
		var v V
		err := binary.Read(buf, binary.LittleEndian, &k)
		if err != nil {
			return nil, err
		}
		err = binary.Read(buf, binary.LittleEndian, &v)
		if err != nil {
			return nil, err
		}
		rr[k] = v
	}

	return rr, nil
}

func MarshalSliceLen[K comparable, V Size](buf io.Writer, rr []K) error {
	// Handle nil slice
	if rr == nil {
		return binary.Write(buf, binary.LittleEndian, maxValue[V]())
	}

	length := V(len(rr))
	// Check for length conflict with nil marker
	if length == maxValue[V]() {
		return errors.New("slice length equals nil marker")
	}

	err := binary.Write(buf, binary.LittleEndian, length)
	if err != nil {
		return err
	}
	err = binary.Write(buf, binary.LittleEndian, rr)
	if err != nil {
		return err
	}
	return nil
}

func UnmarhsalSliceLen[K comparable, V Size](buf io.Reader) ([]K, error) {
	var length V

	err := binary.Read(buf, binary.LittleEndian, &length)
	if err != nil {
		return nil, err
	}

	// Check for nil marker
	if length == maxValue[V]() {
		return nil, nil
	}

	x := make([]K, length)

	err = binary.Read(buf, binary.LittleEndian, &x)
	if err != nil {
		return nil, err
	}
	return x, nil
}

func MarshalWithLen[T Size](buf io.Writer, m encoding.BinaryMarshaler) error {
	// Check if m is nil or is an interface containing a nil pointer
	if m == nil {
		return binary.Write(buf, binary.LittleEndian, maxValue[T]())
	}

	v := reflect.ValueOf(m)
	if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return binary.Write(buf, binary.LittleEndian, maxValue[T]())
		}
	}

	data, err := m.MarshalBinary()
	if err != nil {
		return err
	}

	// Ensure data length doesn't conflict with nil marker
	if T(len(data)) == maxValue[T]() {
		return errors.New("data length equals nil marker")
	}

	length := T(len(data))
	err = binary.Write(buf, binary.LittleEndian, length)
	if err != nil {
		return err
	}

	_, err = buf.Write(data)
	return err
}

func UnmarshalWithLen[T Size](buf io.Reader, m encoding.BinaryUnmarshaler) error {
	var length T

	err := binary.Read(buf, binary.LittleEndian, &length)
	if err != nil {
		return err
	}

	// Check for nil marker
	if length == maxValue[T]() {
		return nil
	}

	data := make([]byte, length)
	_, err = io.ReadFull(buf, data)
	if err != nil {
		return err
	}

	err = m.UnmarshalBinary(data)
	if err != nil {
		return err
	}

	return nil
}

func UnmarshalWithLenNil[T Size](buf io.Reader, m encoding.BinaryUnmarshaler) (hasVal bool, err error) {
	var length T

	err = binary.Read(buf, binary.LittleEndian, &length)
	if err != nil {
		return false, err
	}

	// Check for nil marker
	if length == maxValue[T]() {
		return false, nil
	}

	data := make([]byte, length)

	_, err = io.ReadFull(buf, data)
	if err != nil {
		return false, err
	}

	err = m.UnmarshalBinary(data)
	if err != nil {
		return false, err
	}

	return true, nil
}

func ReadSize[T Size](buf io.Reader) (T, bool, error) {
	var length T

	err := binary.Read(buf, binary.LittleEndian, &length)
	if err != nil {
		return 0, false, err
	}

	return length, length == maxValue[T](), nil
}
