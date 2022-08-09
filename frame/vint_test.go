package frame

import (
	"math"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDecodeVint(t *testing.T) {
	tests := []struct {
		name   string
		data   []byte
		value  int64
		error  bool
		length int
	}{
		{
			name:  "empty",
			data:  []byte{},
			error: true,
		},
		{
			name:   "zero",
			data:   []byte{0x00},
			value:  0,
			length: 1,
		},
		{
			name:   "256000",
			data:   []byte{0xc7, 0xd0, 0x00},
			value:  256000,
			length: 3,
		},
		{
			name:   "123456789",
			data:   []byte{0xee, 0xb7, 0x9a, 0x29},
			value:  -123456789,
			length: 4,
		},
		{
			name:   "MaxInt64",
			data:   []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xfe},
			value:  math.MaxInt64,
			length: 9,
		},
		{
			name:   "MinInt64",
			data:   []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
			value:  math.MinInt64,
			length: 9,
		},
		{
			name:   "more data",
			data:   []byte{0xc7, 0xd0, 0x00, 0x01, 0x02, 0x03},
			value:  256000,
			length: 3,
		},
		{
			name:  "short read",
			data:  []byte{0xff, 0xff, 0xff, 0xff},
			error: true,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			value, length, err := decodeVInt(tc.data)
			if tc.error {
				if err == nil {
					t.Fatal("expected error, but got nil")
				}
				return
			}
			if diff := cmp.Diff(value, tc.value); diff != "" {
				t.Fatalf("value differs\n%s", diff)
			}
			if diff := cmp.Diff(length, tc.length); diff != "" {
				t.Fatalf("length differs\n%s", diff)
			}
		})
	}
}

func TestAppendVint(t *testing.T) {
	tests := []struct {
		name  string
		data  []byte
		value int64
	}{
		{
			name:  "zero",
			data:  []byte{0x00},
			value: 0,
		},
		{
			name:  "256000",
			data:  []byte{0xc7, 0xd0, 0x00},
			value: 256000,
		},
		{
			name:  "123456789",
			data:  []byte{0xee, 0xb7, 0x9a, 0x29},
			value: -123456789,
		},
		{
			name:  "MaxInt64",
			data:  []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xfe},
			value: math.MaxInt64,
		},
		{
			name:  "MinInt64",
			data:  []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
			value: math.MinInt64,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			data := []byte{0x01, 0x02}
			data = appendVInt(data, tc.value)

			expected := []byte{0x01, 0x02}
			expected = append(expected, tc.data...)

			if diff := cmp.Diff(data, expected); diff != "" {
				t.Fatalf("data differs\n%s", diff)
			}
		})
	}
}
