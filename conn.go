package stcp

import (
	"errors"
	"fmt"
	"net"
)

// type conn struct {
// 	Authenticator
// }

type Conn struct {
	conn net.Conn

	rw       *packet
	Session  string
	Sequence int
}

func newConn(tcp net.Conn, proto *packet, cs *ClienState) *Conn {
	sconn := Conn{
		conn:     nil,
		rw:       proto,
		Session:  cs.Session,
		Sequence: cs.Sequence,
	}
	return &sconn
}

func (c *Conn) Conn() net.Conn {
	return c.conn
}

// close connection
func (c *Conn) Close() error {
	return c.conn.Close()
}

// dial to stcp server with default auth
func Dial(addr, username, password string) (*Conn, error) {
	return DialWithSession(addr, username, password, "", 1)
}

// dial to stcp server with session
func DialWithSession(addr, username, password, session string, sequence int) (*Conn, error) {
	tcp, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	//send login request
	rw := NewReaderWriter(tcp, tcp)
	if err := rw.WriteMessage(createLoginRequestMessage(username, password, "", 0)); err != nil {
		return nil, err
	}

	//read login response
	msg, err := rw.ReadMessage()
	if err != nil {
		return nil, err
	}

	//parse login response
	cs := new(ClienState)
	cs, err = loginResponse(msg)
	if err != nil {
		return nil, err
	}

	//connection can use
	c := newConn(tcp, rw, cs)

	return c, nil
}

////////////////////////////////////////////////////

type tcplistener struct {
	l    net.Listener
	conn *Conn
}

func Listen(addr string) (*tcplistener, error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	listener := tcplistener{
		l: l,
	}
	return &listener, nil

}

func (tl *tcplistener) EndSession() error {
	msg := CreateEndSessionMessage()
	err := tl.conn.rw.writer.WriteMessage(msg) //tl.conn.WriteMessage(msg)
	if err != nil {
		return err
	}
	return tl.conn.Close()
}

func (tl *tcplistener) Accept(auth Authenticator) (*Conn, error) {

	tcp, _ := tl.l.Accept()
	soupbin, err := accept(tcp, auth)
	if err != nil {
		return nil, err
	}
	tl.conn = soupbin

	return soupbin, err
}

func accept(conn net.Conn, auth Authenticator) (*Conn, error) {

	rw := NewReaderWriter(conn, conn)
	packet, err := rw.ReadMessage()
	if err != nil {
		return nil, err
	}
	//process auth
	cs, err := clientAuth(packet, auth)
	if err != nil {
		if err == ErrNotAuthorize {
			//create login reject response
			return nil, errors.Join(rw.WriteMessage(createLoginRejectMessage(REJECT_AUTH)), err)
		}
		if err == ErrSessionUnavailabe {
			// create login reject response
			return nil, rw.WriteMessage(createLoginRejectMessage(REJECT_SESSION))
		}
		return nil, err
	}

	err = rw.WriteMessage(createLoginAcceptMessage(cs.Session, cs.Sequence))
	if err != nil {
		return nil, err
	}

	/* can use connection*/
	//heartbeat and send message
	sConn := Conn{
		conn:     conn,
		rw:       rw,
		Session:  cs.Session,
		Sequence: cs.Sequence,
	}
	fmt.Printf("[DEBUG]connect with session %v, sequence %v\n", sConn.Session, sConn.Sequence)

	return &sConn, nil
}
