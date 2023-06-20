package stcp

import (
	"bufio"
	"io"
)

// Writer holds a reference to the underlying writer
// and a buffered writer for efficient writes.
type Writer struct {
	// dst io.Writer     // underlying writer
	bw *bufio.Writer // buffered writer
}

func NewWriter(dst io.Writer) *Writer {
	w := Writer{
		// dst: dst,
		bw: bufio.NewWriter(dst),
	}
	return &w
}

// WriteMessage method takes a message type (typ) and a byte slice (p) as input,
// and it perform the sequencing and writing process.
func (w *Writer) WritePacket(typ uint8, p []byte) error {
	err := w.sequence(typ, p[:])
	if err != nil {
		return err
	}
	return nil
}

// WriteMessage writes the sequenced message with the given type to the underlying writer.
func (w *Writer) sequence(typ uint8, message []byte) error {
	messageLength := len(message)
	packetSize := 1 + messageLength // type byte + p byte as prefix header

	//write prefix [2]byte
	if err := writePrefix(uint16(packetSize), w.bw); err != nil {
		return err
	}

	//write typ uint8/byte
	if err := w.bw.WriteByte(byte(typ)); err != nil {
		return err
	}

	//write message [n]byte
	if _, err := w.bw.Write(message); err != nil {
		return err
	}

	// Flush the buffered writer to the underlying writer
	if err := w.bw.Flush(); err != nil {
		return err
	}

	return nil

}

// writePrefix appends the size bytes to the buffered writer.
func writePrefix(size uint16, buf io.ByteWriter) error {
	if err := buf.WriteByte(byte(size >> 8)); err != nil {
		return err
	}
	if err := buf.WriteByte(byte(size)); err != nil {
		return err
	}

	return nil
}
