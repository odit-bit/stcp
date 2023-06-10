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

func Connect(addr, username, password string, sh SequenceHandlerFunc) (*client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	c := &client{
		session:  "",
		sequence: 0,
		username: &username,
		password: &password,
		sh:       sh,
		pr:       NewReader(conn), //&Reader{buf: [512]byte{}, r: conn},
		pw:       NewWriter(conn), //&Writer{buf: [512]byte{}, w: conn},
		rwc:      conn,
		resetC:   make(chan time.Duration),
	}
	err = c.login(*c.username, *c.password, c.session, c.sequence)
	if err != nil {
		return nil, err
	}
	go c.startHeartbeats()

	return c, nil

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
	fmt.Println("debug login", string(msg))

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
			err := c.sh(b[1:])
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
			return fmt.Errorf("unknown packet %v", string(b[0]))
		}
	}
}

func (c *client) handleLoginAccept(packet []byte) error {

	// parseLoginAccept
	var err error
	c.session = string(bytes.TrimSpace(packet[1:7]))
	c.sequence, err = strconv.Atoi(string(bytes.TrimSpace(packet[7:17])))
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

func (c *client) Close() error {
	close(c.resetC)
	return c.rwc.Close()
}
