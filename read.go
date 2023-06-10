package stcp

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

type Reader struct {
	buf [512]byte
	r   *bufio.Reader // add layer before the underlying reader
}

func NewReader(src io.Reader) *Reader {
	r := Reader{
		buf: [512]byte{},
		r:   bufio.NewReaderSize(src, DefaultBuffer),
	}
	return &r
}

// readMessage read message from a complete  packet,
func (r *Reader) ReadMessage() ([]byte, error) {

	_, err := io.ReadFull(r.r, r.buf[:2])
	if err != nil {
		return nil, err
	}

	size := uint16(r.buf[1]) | uint16(r.buf[0])<<8
	if len(r.buf) < int(size) {
		return nil, fmt.Errorf("not sufficient buffer length for read %v", size)
	}

	_, err = io.ReadFull(r.r, r.buf[0:size])
	if err != nil {
		return nil, err
	}
	return r.buf[0:size], nil
}

// not clear with intention,
func (r *Reader) Read(p []byte) (int, error) {
	// Read from Reader via ReadMessage()
	msg, err := r.ReadMessage()
	if err != nil {
		return 0, err
	}
	buffer := bytes.NewBuffer(msg)
	return buffer.Read(p)
}

// func read(buf []byte, src io.Reader) (int, error) {
// 	var size uint16
// 	if _, err := io.ReadFull(src, buf[:2]); err != nil {
// 		return 0, err
// 	}

// 	size = binary.BigEndian.Uint16(buf[:2])
// 	start := 2 // the first 2 readed byte
// 	n := int(size) + 2

// 	if len(buf) < n {
// 		return n, fmt.Errorf("short buf, need %v", n)
// 	}

// 	// we knew the end (size) of message
// 	// now consume (read) byte until size
// 	if o, err := io.ReadFull(src, buf[start:n]); err != nil {
// 		return 2, err
// 	} else {
// 		n += int(o)
// 	}
// 	return n, nil
// }
