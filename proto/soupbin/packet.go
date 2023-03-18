package soupbin

//packet implemented soupbinTCP binary-protocol
//wrap message as a logical-packet

import (
	"encoding/binary"
	"errors"
	"io"
	"log"
)

const (
	_SESSION_END_TYPE = 'Z'
	_SEQUENCE_TYPE    = 'S'
)

type packetReader struct {
	r   io.Reader
	buf []byte
	//size int
}

// will read binary-packet from underlying reader and unwrap the message as []byte
func (p *packetReader) ReadMessage() ([]byte, error) {

	size, err := p.parseLength()
	if err != nil {
		return nil, err
	}

	n, err := io.ReadFull(p.r, p.buf[0:size])
	if n > size {
		log.Println("out of bound ", n)
		return nil, err
	}

	b := p.buf[0:size]
	p.buf = p.buf[:0]
	return b, nil
}

func (p *packetReader) parseLength() (int, error) {
	//buf := [2]byte{}
	n, err := io.ReadFull(p.r, p.buf[0:2])
	if n != 2 {
		return n, errors.Join(errors.New("wrong size packet"), err)
	}
	length := binary.BigEndian.Uint16(p.buf[0:2])
	return int(length), nil

}

type packetWriter struct {
	w   io.Writer
	buf []byte
	//size int
}

// wrap the message as binary-packet and write to underlying writer
func (p *packetWriter) WriteMessage(msg []byte) error {

	n := p.wrapPacket(msg)
	if n > 0 {
		_, err := p.w.Write(p.buf[0:n])
		if err != nil {
			return err
		}
		p.buf = p.buf[:0]
	}
	return nil

}

// wraping msg as binary-packet into buf. return n bytes was copied to buf
func (p *packetWriter) wrapPacket(msg []byte) int {
	packetLength := len(msg)
	binary.BigEndian.PutUint16(p.buf[0:2], uint16(packetLength))

	// size := packetLength + 2
	n := copy(p.buf[2:packetLength+2], msg)
	return n + 2
}

type packet struct {
	reader *packetReader
	writer *packetWriter
}

func NewReaderWriter(r io.Reader, w io.Writer) *packet {
	p := packet{
		reader: &packetReader{
			r:   r,
			buf: make([]byte, 128),
			//size: 0,
		},
		writer: &packetWriter{
			w:   w,
			buf: make([]byte, 128),
			//size: 0,
		},
	}
	return &p
}

func (p *packet) ReadMessage() ([]byte, error) {
	return p.reader.ReadMessage()
}

func (p *packet) WriteMessage(msg []byte) error {
	return p.writer.WriteMessage(msg)
}
