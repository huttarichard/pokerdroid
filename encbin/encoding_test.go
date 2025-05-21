package encbin

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

type SampleStruct struct {
	A int32
	B float64
}

func (s *SampleStruct) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, s.A); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, s.B); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *SampleStruct) UnmarshalBinary(data []byte) error {
	buf := bytes.NewReader(data)
	if err := binary.Read(buf, binary.LittleEndian, &s.A); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.LittleEndian, &s.B); err != nil {
		return err
	}
	return nil
}

func TestSerializationFunctions(t *testing.T) {
	t.Run("TestMarshalMap", func(t *testing.T) {
		mapData := map[int8]int8{
			3: 1,
			4: 2,
			5: 3,
		}
		buf := new(bytes.Buffer)
		err := MarshalMap[int8, int8, uint16](buf, mapData)
		require.NoError(t, err)

		outMap, err := UnmarshalMap[int8, int8, uint16](buf)
		require.NoError(t, err)
		require.Equal(t, mapData, outMap)
	})

	t.Run("TestMarshalSliceLen", func(t *testing.T) {
		sliceData := []int8{1, 2, 3}
		buf := new(bytes.Buffer)
		err := MarshalSliceLen[int8, uint16](buf, sliceData)
		require.NoError(t, err)

		outSlice, err := UnmarhsalSliceLen[int8, uint16](buf)
		require.NoError(t, err)
		require.Equal(t, sliceData, outSlice)
	})

	t.Run("TestMarshalValues", func(t *testing.T) {
		values := []interface{}{int32(5), float64(7.5)}
		buf := new(bytes.Buffer)
		err := MarshalValues(buf, values...)
		require.NoError(t, err)

		var i int32
		var f float64
		err = UnmarshalValues(buf, &i, &f)
		require.NoError(t, err)
		require.Equal(t, int32(5), i)
		require.Equal(t, 7.5, f)
	})

	t.Run("TestMarshalWithLen", func(t *testing.T) {
		data := &SampleStruct{
			A: 42,
			B: 42.42,
		}
		buf := new(bytes.Buffer)
		err := MarshalWithLen[uint16](buf, data)
		require.NoError(t, err)

		outData := &SampleStruct{}
		err = UnmarshalWithLen[uint16](buf, outData)
		require.NoError(t, err)
		require.Equal(t, data, outData)
	})

	t.Run("TestMarshalWithLenNilHandling", func(t *testing.T) {
		t.Run("max value handling", func(t *testing.T) {
			// Create data that will marshal to maxUint8 bytes
			s := &testMarshaler{
				Data: make([]byte, math.MaxUint8),
			}

			buf := new(bytes.Buffer)
			err := MarshalWithLen[uint8](buf, s)
			require.Error(t, err)
			require.Contains(t, err.Error(), "data length equals nil marker")
		})

		t.Run("nil interface", func(t *testing.T) {
			var m encoding.BinaryMarshaler
			buf := new(bytes.Buffer)

			// Marshal nil
			err := MarshalWithLen[uint8](buf, m)
			require.NoError(t, err)

			// Unmarshal
			var out SampleStruct
			err = UnmarshalWithLen[uint8](buf, &out)
			require.NoError(t, err)
			require.Equal(t, SampleStruct{}, out)
		})

		t.Run("nil pointer", func(t *testing.T) {
			var m *SampleStruct
			buf := new(bytes.Buffer)

			// Marshal nil pointer
			err := MarshalWithLen[uint16](buf, m)
			require.NoError(t, err)

			// Unmarshal
			out := &SampleStruct{}
			err = UnmarshalWithLen[uint16](buf, out)
			require.NoError(t, err)
			require.Equal(t, &SampleStruct{}, out)
		})
	})
}

func TestNilHandling(t *testing.T) {
	t.Run("nil map handling", func(t *testing.T) {
		t.Run("marshal nil map", func(t *testing.T) {
			var m map[int8]int8
			buf := new(bytes.Buffer)
			err := MarshalMap[int8, int8, uint8](buf, m)
			require.NoError(t, err)

			outMap, err := UnmarshalMap[int8, int8, uint8](buf)
			require.NoError(t, err)
			require.Nil(t, outMap)
		})

		t.Run("large map", func(t *testing.T) {
			// Create map with maxUint8 entries
			m := make(map[int8]int8)
			for i := int8(0); i < math.MaxInt8; i++ {
				m[i] = i
			}
			// Add enough entries to reach maxUint8
			for i := int8(0); len(m) < math.MaxUint8; i++ {
				m[i] = i
			}

			buf := new(bytes.Buffer)
			err := MarshalMap[int8, int8, uint8](buf, m)
			require.Error(t, err)
			require.Contains(t, err.Error(), "map length equals nil marker")
		})
	})

	t.Run("nil slice handling", func(t *testing.T) {
		t.Run("marshal nil slice", func(t *testing.T) {
			var s []int8
			buf := new(bytes.Buffer)
			err := MarshalSliceLen[int8, uint8](buf, s)
			require.NoError(t, err)

			outSlice, err := UnmarhsalSliceLen[int8, uint8](buf)
			require.NoError(t, err)
			require.Nil(t, outSlice)
		})

		t.Run("large slice", func(t *testing.T) {
			// Create slice with maxUint8 entries
			s := make([]int8, math.MaxUint8)

			buf := new(bytes.Buffer)
			err := MarshalSliceLen[int8, uint8](buf, s)
			require.Error(t, err)
			require.Contains(t, err.Error(), "slice length equals nil marker")
		})
	})
}

// Add this helper type that implements BinaryMarshaler
type testMarshaler struct {
	Data []byte
}

func (t *testMarshaler) MarshalBinary() ([]byte, error) {
	return t.Data, nil
}

func (t *testMarshaler) UnmarshalBinary(data []byte) error {
	t.Data = data
	return nil
}

// Test helper types
type testStruct struct {
	A int32
	B float64
}

func (s *testStruct) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := MarshalValues(buf, s.A, s.B)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *testStruct) UnmarshalBinary(data []byte) error {
	return UnmarshalValues(bytes.NewReader(data), &s.A, &s.B)
}

func TestMarshalWithLen(t *testing.T) {
	tests := []struct {
		name    string
		value   encoding.BinaryMarshaler
		wantErr bool
	}{
		{
			name:    "nil value",
			value:   nil,
			wantErr: false,
		},
		{
			name:    "basic struct",
			value:   &testStruct{A: 42, B: 3.14},
			wantErr: false,
		},
		{
			name: "max length data",
			value: &testStruct{
				A: 1<<31 - 1,
				B: 1.797693134862315708145274237317043567981e+308,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal
			buf := new(bytes.Buffer)
			err := MarshalWithLen[uint32](buf, tt.value)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Unmarshal
			var out testStruct
			err = UnmarshalWithLen[uint32](buf, &out)
			require.NoError(t, err)

			// Compare
			if tt.value == nil {
				return
			}
			require.Equal(t, tt.value, &out)
		})
	}
}

func TestMarshalMap(t *testing.T) {
	tests := []struct {
		name    string
		value   map[int32]float64
		wantErr bool
	}{
		{
			name:    "nil map",
			value:   nil,
			wantErr: false,
		},
		{
			name:    "empty map",
			value:   map[int32]float64{},
			wantErr: false,
		},
		{
			name: "basic map",
			value: map[int32]float64{
				1: 2.5,
				2: 3.14,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal
			buf := new(bytes.Buffer)
			err := MarshalMap[int32, float64, uint16](buf, tt.value)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Unmarshal
			out, err := UnmarshalMap[int32, float64, uint16](buf)
			require.NoError(t, err)

			// Compare
			require.Equal(t, tt.value, out)
		})
	}
}

func TestMarshalSliceLen(t *testing.T) {
	tests := []struct {
		name    string
		value   []int32
		wantErr bool
	}{
		{
			name:    "nil slice",
			value:   nil,
			wantErr: false,
		},
		{
			name:    "empty slice",
			value:   []int32{},
			wantErr: false,
		},
		{
			name:    "basic slice",
			value:   []int32{1, 2, 3},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal
			buf := new(bytes.Buffer)
			err := MarshalSliceLen[int32, uint16](buf, tt.value)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Unmarshal
			out, err := UnmarhsalSliceLen[int32, uint16](buf)
			require.NoError(t, err)

			// Compare
			require.Equal(t, tt.value, out)
		})
	}
}

func TestUnmarshalWithLenNil(t *testing.T) {
	tests := []struct {
		name       string
		value      encoding.BinaryMarshaler
		wantHasVal bool
		wantErr    bool
	}{
		{
			name:       "nil value",
			value:      nil,
			wantHasVal: false,
			wantErr:    false,
		},
		{
			name:       "basic struct",
			value:      &testStruct{A: 42, B: 3.14},
			wantHasVal: true,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal
			buf := new(bytes.Buffer)
			err := MarshalWithLen[uint32](buf, tt.value)
			require.NoError(t, err)

			// Unmarshal
			var out testStruct
			hasVal, err := UnmarshalWithLenNil[uint32](buf, &out)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.wantHasVal, hasVal)

			if hasVal {
				require.Equal(t, tt.value, &out)
			}
		})
	}
}
