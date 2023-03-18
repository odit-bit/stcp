package stcp

import (
	"bytes"
	"testing"
)

func Benchmark_createLoginRequestMsg(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b := createLoginRequestMessage("admin", "12345", "22a10", 852)
		_ = b
	}
}

func Test_packetReader(t *testing.T) {
	data := []byte{0, 9, 65, 32, 32, 32, 32, 32, 32, 65, 65}
	buf := bytes.NewBuffer(data)
	p := NewReaderWriter(buf, buf)

	actual, err := p.reader.ReadMessage()
	if err != nil {
		t.Errorf("want nil got %v", err)
	}
	expected := []byte{65, 32, 32, 32, 32, 32, 32, 65, 65}
	for i, v := range expected {
		if v != actual[i] {
			t.Errorf("want %v, got %v ", v, actual[i])
		}
	}
}

func Test_packetWriter(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	p := NewReaderWriter(buf, buf)

	data := []byte{65, 32, 32, 32, 32, 32, 32, 65, 65}
	err := p.writer.WriteMessage(data)
	if err != nil {
		t.Errorf("want nil got %v", err)
	}
	actual := buf.Bytes()
	expected := []byte{0, 9, 65, 32, 32, 32, 32, 32, 32, 65, 65}
	if len(actual) != len(expected) {
		t.Fatalf("want %v, got %v", len(expected), len(actual))
	}
	for i, v := range expected {
		if v != actual[i] {
			t.Errorf("want %v, got %v ", v, actual[i])
		}
	}
}

func Benchmark_PacketReader(b *testing.B) {
	data := []byte{0, 9, 65, 32, 32, 32, 32, 32, 32, 65, 65}
	buf := bytes.NewReader(data)
	p := NewReaderWriter(buf, bytes.NewBuffer([]byte{}))

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		buf.Reset(data)
		msg, err := p.reader.ReadMessage()
		if err != nil {
			b.Fatal(err)
		}

		if len(msg) != 9 {
			b.Fatal("length not 9")
		}
		if msg[8] != 65 {
			b.Fatalf("element not 65 got %v", msg[8])
		}
	}
}

func Benchmark_packetWriter(b *testing.B) {
	dst := []byte{} // stand-in as destination buffer
	buf := bytes.NewBuffer(dst)

	p := NewReaderWriter(buf, buf)
	msg := []byte{65, 32, 32, 32, 32, 32, 32, 65, 65}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		err := p.writer.WriteMessage(msg)
		if err != nil {
			b.Fatal(err)
		}
		buf.Reset()
	}
}
