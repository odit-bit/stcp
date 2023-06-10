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

func (w *Writer) WriteMessage(typ uint8, p []byte) error {
	err := w.sequence(typ, p[:])
	if err != nil {
		return err
	}
	return nil
}

// sequence sequencing the message as packet befor write to underlying writer
func (w *Writer) sequence(typ uint8, payload []byte) error {
	n := len(payload)
	size := 1 + n // type byte + p byte as prefix header
	if err := writePrefix(typ, uint16(size), w.buf); err != nil {
		return err
	}

	if _, err := w.buf.Write(payload); err != nil {
		return err
	}
	if _, err := w.dst.Write(w.buf.Bytes()); err != nil {
		return err
	}

	return nil
}

// implement to write prefix ( size and type ) before write the actual payload
func writePrefix(typ uint8, size uint16, buf *bytes.Buffer) error {
	if err := buf.WriteByte(byte(size >> 8)); err != nil {
		return err
	}
	if err := buf.WriteByte(byte(size)); err != nil {
		return err
	}

	if err := buf.WriteByte(byte(typ)); err != nil {
		return err
	}
	return nil
}

//////////////////////////
