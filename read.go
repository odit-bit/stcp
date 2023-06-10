package stcp

import (
	"bufio"
	"fmt"
	"io"
)

type Reader struct {
	buf [512]byte
	r   *bufio.Reader //  add layer before the underlying reader
}

func NewReader(src io.Reader) *Reader {
	r := Reader{
		buf: [512]byte{},
		r:   bufio.NewReaderSize(src, DefaultBuffer),
	}
	return &r
}

// ReadMessage return a complete sequence packet byte, removed the prefix header
func (r *Reader) ReadMessage() ([]byte, error) {

	_, err := io.ReadFull(r.r, r.buf[:2])
	if err != nil {
		return nil, err
	}

	payloadLength := uint16(r.buf[1]) | uint16(r.buf[0])<<8
	if len(r.buf) < int(payloadLength) {
		return nil, fmt.Errorf("not sufficient buffer length for read %v", payloadLength)
	}

	end := payloadLength
	_, err = io.ReadFull(r.r, r.buf[0:end])
	if err != nil {
		return nil, err
	}

	return r.buf[0:end], nil
}
