package stcp

import (
	"bytes"
	"io"
)

// Writer that can provide  buffering and sequencing message before write it to underlying writer
type Writer struct {
	dst io.Writer     // underlying writer
	buf *bytes.Buffer // internal buffer
}

func NewWriter(dst io.Writer) *Writer {
	w := Writer{
		dst: dst,
		buf: &bytes.Buffer{}, //[512]byte{},
	}
	return &w
}

// WriteMessage packed p to network type and write to underlying writer
func (w *Writer) WriteMessage(typ uint8, p []byte) error {
	err := w.sequence(typ, p[:])
	if err != nil {
		return err
	}
	return nil
}

// actual implementation of sequencing the message as packet befor write to underlying writer
func (w *Writer) sequence(typ uint8, message []byte) error {
	messageLength := len(message)
	packetSize := 1 + messageLength // type byte + p byte as prefix header

	//write prefix [2]byte
	if err := writePrefix(uint16(packetSize), w.buf); err != nil {
		return err
	}

	//write typ uint8/byte
	if err := w.buf.WriteByte(byte(typ)); err != nil {
		return err
	}

	//write message [n]byte
	if _, err := w.buf.Write(message); err != nil {
		return err
	}

	//actual read to underlyng writer
	if _, err := w.dst.Write(w.buf.Bytes()); err != nil {
		return err
	}
	w.buf.Reset()
	return nil
}

// writePrefix append (size) byte to buf
func writePrefix(size uint16, buf *bytes.Buffer) error {
	if err := buf.WriteByte(byte(size >> 8)); err != nil {
		return err
	}
	if err := buf.WriteByte(byte(size)); err != nil {
		return err
	}

	return nil
}
