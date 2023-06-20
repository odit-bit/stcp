package stcp

import (
	"bytes"
	"io"
	"testing"
)

func TestWrite(t *testing.T) {
	var tt = []struct {
		name   string
		input  []byte
		typ    uint8
		expect []byte
	}{
		{
			name:   "Write Sequence packet",
			input:  []byte("hello"),
			typ:    'S',
			expect: []byte{0, 6, 'S', 'h', 'e', 'l', 'l', 'o'},
		},
	}

	var dst bytes.Buffer

	w := NewWriter(&dst)
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := w.WritePacket(tc.typ, tc.input[:])
			if err != nil {
				t.Error(err)
			}

			actual := dst.Bytes()
			if !bytes.Equal(actual, tc.expect) {
				t.Errorf("got %v", actual)
			}
			dst.Reset()
		})
	}
}

func Benchmark_Write(b *testing.B) {
	input := []byte{'h', 'e', 'l', 'l', 'o'}
	typ := uint8('S')
	var dst = io.Discard

	w := NewWriter(dst)

	var result []byte

	loopCount := b.N
	b.ResetTimer()
	for i := 0; i < loopCount; i++ {
		err := w.WritePacket(typ, input)
		if err != nil {
			b.Fatal(err)
		}

	}
	BencResult = result
}
