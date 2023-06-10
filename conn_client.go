package stcp

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"
)

type SequenceHandlerFunc func([]byte) error

type ClientConfig struct {
	Addr     string
	Username string
	Password string

	SH SequenceHandlerFunc
}

type client struct {
	session  string
	sequence int
	username *string
	password *string

	sh  SequenceHandlerFunc
	pr  SequenceReader
	pw  SequenceWriter
	rwc net.Conn // rwc Transport

	resetC chan time.Duration
}

func NewClient(conf *ClientConfig) (*client, error) {
	conn, err := net.Dial("tcp", conf.Addr)
	if err != nil {
		return nil, err
	}

	c := &client{
		session:  "",
		sequence: 0,
		username: &conf.Username,
		password: &conf.Password,
		sh:       conf.SH,
		pr:       NewReader(conn), //&Reader{buf: [512]byte{}, r: conn},
		pw:       NewWriter(conn), //&Writer{buf: [512]byte{}, w: conn},
		rwc:      conn,
		resetC:   make(chan time.Duration),
	}
	return c, nil
}

func NewClienAndConnect(conf *ClientConfig) error {
	c, err := NewClient(conf)
	if err != nil {
		return err
	}
	if err := c.Connect(); err != nil {
		return err
	}
	defer func() {
		fmt.Println("Sequence packet received", c.SequenceNumber())
		c.rwc.Close()
	}()

	if err := c.Receive(); err != nil {
		return err
	}
	return nil

}

func (c *client) Connect() error {

	if err := c.login(*c.username, *c.password, c.session, c.sequence); err != nil {
		return err
	}

	go c.startHeartbeats()

	return nil
}

func (c *client) sendLoginRequest(username, password, session, sequence string) error {
	lr := []byte{}
	lr = append(lr, []byte(fmt.Sprintf("%-6s", username))...)
	lr = append(lr, []byte(fmt.Sprintf("%-10s", password))...)
	lr = append(lr, []byte(fmt.Sprintf("%-10s", session))...)
	lr = append(lr, []byte(fmt.Sprintf("%-20s", sequence))...)

	if err := c.pw.WriteMessage('L', lr); err != nil {
		return err
	}
	return nil
}

func (c *client) login(username, password, session string, sequence int) error {
	//send login request
	if err := c.sendLoginRequest(username, password, session, strconv.Itoa(sequence)); err != nil {
		return err
	}

	//retreive login response
	msg, err := c.pr.ReadMessage()
	if err != nil {
		return err
	}
	switch msg[0] {
	case LoginAcceptType:
		//loginAccept
		if err := c.handleLoginAccept(msg); err != nil {
			return err
		}
	case LoginRejectType:
		//loginReject
		return fmt.Errorf("login rejected with reason code %v", string(msg[1]))
	}
	return nil
}

// receive packet and invoke handler according to packet type
func (c *client) Receive() error {
	if c.sh == nil {
		log.Println("using default handler for Sequence Packet")
		c.sh = func(b []byte) error {
			fmt.Printf("sequence length :%v payload : %x \n", string(b[0]), b[3:])
			return nil
		}
	}

	for {
		c.rwc.SetReadDeadline(time.Now().Add(ReadTimeoutDefault))
		b, err := c.pr.ReadMessage()
		if err != nil {
			return err
		}
		switch b[0] {
		case SequenceMessage:
			err := c.sh(b)
			if err != nil {
				return err
			}
			c.sequence++
		case 'Z':
			log.Println("end of session")
			return nil
		case ServerHbType:
			continue
		default:
			return fmt.Errorf("unknown packet")
		}
	}
}

func (c *client) handleLoginAccept(packet []byte) error {
	msg, err := c.pr.ReadMessage()
	if err != nil {
		return err
	}

	// parseLoginAccept
	c.session = string(bytes.TrimSpace(msg[1:7]))
	c.sequence, err = strconv.Atoi(string(bytes.TrimSpace(msg[7:17])))
	if err != nil {
		return err
	}

	return nil
}

func (c *client) startHeartbeats() {
	heartbeats(ClientHbType, c.resetC, c.rwc)
}

func (c *client) SequenceNumber() int {
	return c.sequence
}
