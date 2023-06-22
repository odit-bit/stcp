package stcp

import (
	"fmt"
	"net"
	"time"
)

// will implemented by object that can Write in form of soupbin wire type (packet)
type SequenceWriter interface {
	WriteMessage(typ uint8, p []byte) error
}

// will implemented by object that can read soupbin wire type (packet)
type SequenceReader interface {
	ReadMessage() ([]byte, error)
}

type LoginValidator interface {
	IsValid(username, password string) bool
	SessionValid(session string, sequence string) (string, string, bool)
}

///////////////////////////////////

// encapsulate soupbin Read and Write operation
type Conn struct {
	conn   net.Conn
	reader *Reader
	writer *Writer

	resetT chan time.Duration
}

func NewConn(rwc net.Conn) *Conn {
	c := Conn{
		conn:   rwc,
		reader: NewReader(rwc),
		writer: NewWriter(rwc),
		resetT: make(chan time.Duration),
	}

	return &c
}

// ReadMessage return complete packet bytes.
// If the peer does not send anything (neither data nor heartbeats) for an
// extended period of time, it can assume that the link is down
func (c *Conn) ReadMessage() ([]byte, error) {
	c.conn.SetReadDeadline(time.Now().Add(TimeoutDefault))
	msg, err := c.reader.ReadMessage()
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (c *Conn) WriteMessage(typ uint8, buffer []byte) error {
	// c.conn.SetWriteDeadline(time.Now().Add(2 * TimeoutDefault))
	err := c.writer.WritePacket(typ, buffer)
	if err != nil {
		return err
	}

	// mitigation if the c.resetT not ready receive the value
	// ex : no go routine yet consume the channel
	select {
	case c.resetT <- HeartbeatTimeout:
	default:
	}

	return nil
}

// it is mechanism to send heartbeat packet while no packet to send ,
// ensures that the peer will receive data on a regular basis.
func (c *Conn) StartHeartbeats(typ uint8) {
	heartbeatsTicker(c, typ, c.resetT)
}

func (c *Conn) Close() error {
	close(c.resetT)
	c.reader = nil
	c.writer = nil

	err := c.conn.Close()
	if err != nil {
		return err
	}
	return nil
}

// implement heartbeats using time.Ticker
func heartbeatsTicker(c *Conn, typ uint8, resetC chan time.Duration) {
	defer func() {
		fmt.Println("exit heartbeat")
	}()

	var interval time.Duration = HeartbeatTimeout

	timer := time.NewTicker(interval)
	data := []byte{}
	for {

		select {
		case _, ok := <-timer.C:
			if !ok {
				return
			}
			err := c.writer.WritePacket(typ, data)
			if err != nil {
				fmt.Println("heartbeat write error ", err)
				return
			}
		case newInterval, ok := <-resetC:
			if !ok {
				timer.Stop()
				return
			}
			timer.Reset(newInterval)
		}
	}
}

// old deprecated  ///
