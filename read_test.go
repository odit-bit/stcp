package stcp

import (
	"bytes"
	"io"
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
	reader := NewReader(&src)

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

func Test_ReaderWriteTo(t *testing.T) {
	var tt = []struct {
		name   string
		input  []byte
		expect []byte       // expected byte readed or Writed to
		prefix int          // this field contain length of payload
		dst    bytes.Buffer // stand in as a destionation writer
	}{
		{
			name:   "writeTo Sequence Packet",
			input:  []byte{0, 6, 'S', 'h', 'e', 'l', 'l', 'o'},
			expect: []byte("Shello"),
			prefix: 8,
			dst:    bytes.Buffer{},
		},
	}

	var src bytes.Reader
	reader := NewReader(&src)

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			src.Reset(tc.input)
			n, err := reader.WriteTo(&tc.dst)
			if err != nil {
				t.Error(err)
			}
			if n != int64(tc.prefix) {
				t.Error(n)
			}
			actual := tc.dst.Bytes()
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

	r := NewReader(&src)

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

func Benchmark_WriteTo(b *testing.B) {
	input := []byte{0, 6, 'S', 'h', 'e', 'l', 'l', 'o'}
	var src bytes.Reader
	var dst = io.Discard //bytes.Buffer
	var result []byte

	r := NewReader(&src)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		src.Reset(input)

		n, err := r.WriteTo(dst)
		if err != nil {
			b.Fatal(err, n)
		}

	}
	BencResult = result
}
