package stcp

import (
	"bufio"
	"bytes"
	"testing"
)

func TestRead(t *testing.T) {
	var tt = []struct {
		name   string
		input  []byte
		expect []byte
	}{
		{
			name:   "Read Sequence Packet",
			input:  []byte{0, 6, 'S', 'h', 'e', 'l', 'l', 'o'},
			expect: []byte("Shello"),
		},
	}

	var src bytes.Reader
	reader := Reader{
		buf: [512]byte{},
		r:   bufio.NewReader(&src),
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			src.Reset(tc.input)
			actual, err := reader.ReadMessage()
			if err != nil {
				t.Error(err)
			}
			if !bytes.Equal(actual, tc.expect) {
				t.Errorf("got %v", string(actual))
			}
		})
	}
}

var BencResult []byte

func Benchmark_Read(b *testing.B) {
	input := []byte{0, 6, 'S', 'h', 'e', 'l', 'l', 'o'}
	var src bytes.Reader

	r := Reader{
		buf: [512]byte{},
		r:   bufio.NewReader(&src),
	}

	var result []byte
	for i := 0; i < b.N; i++ {
		src.Reset(input)
		res, err := r.ReadMessage()
		if err != nil {
			b.Fatal(err)
		}

		result = res
	}
	BencResult = result
}
