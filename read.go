package stcp

import (
	"bufio"
	"fmt"
	"io"
	"time"
)

const defaultReadBuffer = 512

type Reader struct {
	prefixBuf  [2]byte
	payloadBuf [512]byte
	r          *bufio.Reader //  add layer before the underlying reader
}

func NewReader(src io.Reader) *Reader {
	r := Reader{
		prefixBuf:  [2]byte{},
		payloadBuf: [512]byte{},
		r:          bufio.NewReaderSize(src, defaultReadBuffer),
	}
	return &r
}

// ReadMessage return a complete sequence packet byte,
// removed the prefix header but retain the type field
func (r *Reader) ReadMessage() ([]byte, error) {

	_, err := io.ReadFull(r.r, r.prefixBuf[:2])
	if err != nil {
		return nil, err
	}

	payloadLength := uint16(r.prefixBuf[1]) | uint16(r.prefixBuf[0])<<8
	if len(r.payloadBuf) < int(payloadLength) {
		return nil, fmt.Errorf("not sufficient buffer length for read %v", r.prefixBuf)
	}

	// Read the payload in chunks until all bytes are read
	// readCount := 0
	// for readCount < int(payloadLength) {
	// 	n, err := r.r.Read(r.payloadBuf[readCount:])
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	readCount += n
	// }
	_, err = io.ReadFull(r.r, r.payloadBuf[:payloadLength])
	if err != nil {
		return nil, err
	}
	return r.payloadBuf[:payloadLength], nil
}

// implement io.WriterTo
func (r *Reader) WriteTo(dst io.Writer) (int64, error) {
	n, err := io.ReadFull(r.r, r.prefixBuf[:2])
	if err != nil {
		return 0, err
	}

	payloadLength := uint16(r.prefixBuf[1]) | uint16(r.prefixBuf[0])<<8
	if len(r.payloadBuf) < int(payloadLength) {
		return int64(n), fmt.Errorf("not sufficient buffer length for read %v", payloadLength)
	}

	o, err := io.CopyN(dst, r.r, int64(payloadLength))
	if err != nil {
		return int64(n), err
	}

	n += int(o)
	return int64(n), nil
}

// it is still Beta.
// ReadMessageStream reads message packets continuously from the underlying reader in a streaming fashion.
// It returns a channel of message payloads that are received asynchronously.
// The function stops reading and closes the channel when an error occurs.
func (r *Reader) ReadMessageStream() <-chan []byte {
	payloads := make(chan []byte)

	go func() {
		defer close(payloads)

		for {
			payload, err := r.ReadMessage()
			if err != nil {
				// An error occurred, stop reading and exit the goroutine
				return
			}

			select {
			case payloads <- payload:
				// Payload sent successfully, continue reading
			case <-time.After(time.Second): // Adjust the timeout duration as needed
				// Timed out waiting to send the payload, continue reading
			}
		}
	}()

	return payloads
}
