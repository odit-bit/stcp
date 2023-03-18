package stcp

import (
	"errors"
	"log"
	"net"
	"time"
)

type Client struct {
	username *string
	password *string
	session  *string
	sequence *int
	address  *string

	seqMsgHandler SeqDataHandlerFunc
	netconn       net.Conn
	Reader        *packetReader
	Writer        *packetWriter

	seqMsgC chan []byte
	stop    chan struct{}
	errChan chan error
}

// inisiate client with default options
func NewClient(username, password string) *Client {
	opts := NewClientOptions()
	opts.SetPassword(password)
	opts.SetUsername(username)
	opts.SetSession("")
	opts.SetSequence(1)
	opts.Address(":6666")

	return opts.client
}

func NewClientWithOptions(opt *ClientOptions) *Client {
	return opt.client
}

func check(c *Client) error {
	if c.seqMsgHandler == nil {
		return errors.New("no EventHandler")
	}
	if c.username == nil {
		return errors.New("username not set")
	}
	if c.password == nil {
		return errors.New("password not set")
	}
	if c.password == nil {
		return errors.New("session not set")
	}
	if c.password == nil {
		return errors.New("sequence not set")
	}
	return nil
}

func (c *Client) Connect() error {
	conn, err := net.Dial("tcp", *c.address)
	if err != nil {
		return err
	}

	c.netconn = conn
	c.Reader = NewReader(conn)
	c.Writer = NewWriter(conn)

	//check necessary evil
	if err := check(c); err != nil {
		return err
	}

	//sent login request to authentication, if auth failed handleEvent will return error
	lr := createLoginRequestMessage(*c.username, *c.password, *c.session, *c.sequence)
	c.Writer.WriteMessage(lr)

	return c.handleEvent()
}
func (c *Client) Close() {
	log.Println("connection close")
	c.netconn.Close()
}

// all channel goroutin will end up here. it can say as endpoint
func (c *Client) handleEvent() error {
	//1
	if c.seqMsgHandler == nil {
		return errors.New("no handler for received sequenced message")
	}

	//2 dispatch reader concurently to read from inbound connection
	//and select handling and channels according to the appropriate type
	go c.readEventWorker(c.Reader, c.seqMsgC, c.errChan)

	//3 all channel goroutine will end up here. it can say as the endpoint
	//escpecially the sequenced message for further processing.
	//while no msg from goroutine for periode of time it will sent heartbeat message,
	//as an indication that the connection is still a live
	hb := CreateHeartBeat('R')
	hbTicker := time.NewTicker(1000 * time.Millisecond)
	for {
		select {
		case msg := <-c.seqMsgC:
			if err := c.handleSeqMsg(msg); err != nil {
				return err
			}

		case <-hbTicker.C:
			if err := c.Writer.WriteMessage(hb); err != nil {
				return err
			}

		case err := <-c.errChan:
			return err

		case <-c.stop:
			return nil
		}
	}

}

// read from inbound connection and Select the appropriate handling.
func (c *Client) readEventWorker(reader ProtoReader, msg chan<- []byte, errC chan<- error) {
	for {
		b, err := reader.ReadMessage()
		if err != nil {
			errC <- err
			return
		}

		t := b[0]
		switch t {
		case _HEARTBEAT_SERVER_TYPE:
			continue

		case _SEQUENCE_TYPE:
			msg <- b

		case _LOGIN_ACCEPT_TYPE:
			if err := c.handleLoginAccept(b); err != nil {
				errC <- err
			}

		//stop case or internal event
		case _LOGIN_REJECT_TYPE:
			errC <- parseLoginRejectMessage(b)
			return

		case _SESSION_END_TYPE:
			log.Println("end session invoke")
			c.stop <- struct{}{}
			return

		default:
			errC <- ErrUnknownPacket
			return
		}
	}
}

func (c *Client) handleSeqMsg(msg []byte) error {
	err := c.seqMsgHandler(msg)
	if err != nil {
		return err
	}
	//*c.sequence += 1
	return nil
}

func (c *Client) handleLoginAccept(msg []byte) error {
	cs, err := parseLoginAccept(msg)
	if err != nil {
		return err
	}
	c.sequence = &cs.Sequence
	c.session = &cs.Session
	return nil
}

// /////
type ClientOptions struct {
	client *Client
}

func NewClientOptions() *ClientOptions {
	c := &Client{
		username:      nil,
		password:      nil,
		session:       new(string),
		sequence:      new(int),
		address:       nil,
		netconn:       nil,
		Reader:        nil,
		Writer:        nil,
		seqMsgHandler: nil,
		seqMsgC:       make(chan []byte, 1),
		stop:          make(chan struct{}, 1),
		errChan:       make(chan error, 1),
	}
	return &ClientOptions{
		client: c,
	}
}

func (opt *ClientOptions) SetUsername(un string) {
	opt.client.username = &un
}

func (opt *ClientOptions) SetPassword(pass string) {
	opt.client.password = &pass
}

func (opt *ClientOptions) SetSession(sess string) {
	opt.client.session = &sess
}

func (opt *ClientOptions) SetSequence(seq int) {
	opt.client.sequence = &seq
}

func (opt *ClientOptions) DataHandlerFunc(fn SeqDataHandlerFunc) {
	opt.client.seqMsgHandler = fn
}

func (opt *ClientOptions) Address(addr string) {
	opt.client.address = &addr
}
