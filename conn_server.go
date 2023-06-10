package stcp

import (
	"fmt"
	"io"
	"net"
	"time"
)

type Authenticator func(username, password, session, sequence string) (string, string, error)

type SequenceWriter interface {
	WriteMessage(typ uint8, p []byte) error
}
type SequenceReader interface {
	ReadMessage() ([]byte, error)
}

// encapsulate server side operation
type serverConn struct {
	rwc net.Conn
	sw  SequenceWriter
	sr  SequenceReader

	authFunc Authenticator
	resetC   chan time.Duration
}

func NewServerWithOpts(conn net.Conn, opts ...ServerOptionFunc) (*serverConn, error) {
	s := &serverConn{
		rwc:    conn,
		sw:     NewWriter(conn), //&Writer{buf: [512]byte{}, w: conn},
		sr:     NewReader(conn), //&Reader{buf: [512]byte{}, r: conn},
		resetC: make(chan time.Duration),
	}

	for _, opt := range opts {
		opt(s)
	}

	if s.authFunc == nil {
		fmt.Println("using default auth handler")
		s.authFunc = func(username, password, session, sequence string) (string, string, error) {
			return "default", "1", nil
		}
	}

	return s, nil
}

func NewServerConn(conn net.Conn, auth Authenticator) *serverConn {
	s := &serverConn{
		rwc:      conn,
		sw:       NewWriter(conn),
		sr:       NewReader(conn),
		authFunc: auth,
		resetC:   make(chan time.Duration),
	}

	return s
}

func Accept(l net.Listener, auth Authenticator) (*serverConn, error) {
	c, err := l.Accept()
	if err != nil {
		return nil, err
	}

	opt := SetAuthandler(auth)
	srv, _ := NewServerWithOpts(c, opt)
	if err := srv.Auth(); err != nil {
		return nil, err
	}

	go heartbeats(ServerHbType, srv.resetC, srv.rwc)
	return srv, nil
}

func (s *serverConn) Auth() error {

	//expect loginrequest
	b, err := s.sr.ReadMessage()
	if err != nil {
		return err
	}

	//decode b and it must loginRequest
	data, err := LoginRequestByteToSlice(b)
	if err != nil {
		return err
	}

	//auth client by data
	session, sequence, err := s.authFunc(data[0], data[1], data[2], data[3])
	if err != nil {
		//send reject
		s.sendLoginReject('A')
		if err != nil {
			return err
		}
		return fmt.Errorf("reject login request with code %v , ", string('A'))
	}

	// send accept
	if err := s.sendLoginAccept(session, sequence); err != nil {
		return err
	}

	go heartbeats(ServerHbType, s.resetC, s.rwc)
	return nil
}

// send payload as sequence message
func (s *serverConn) Send(payload []byte) error {
	err := s.sw.WriteMessage('S', payload)
	if err != nil {
		return err
	}
	s.resetC <- HeartbeatTimeout
	return nil
}

func (s *serverConn) sendLoginReject(reason uint8) error {
	return s.sw.WriteMessage(LoginRejectType, []byte{reason})
}

func (s *serverConn) sendLoginAccept(session, sequence string) error {
	msg := []byte{}
	msg = append(msg, []byte(fmt.Sprintf("%-10s", session))...)
	msg = append(msg, []byte(fmt.Sprintf("%-20s", sequence))...)

	return s.sw.WriteMessage(LoginAcceptType, msg)
}

// receive message byte
func (s *serverConn) Receive() ([]byte, error) {
	for {
		s.rwc.SetReadDeadline(time.Now().Add(ReadTimeoutDefault))
		b, err := s.sr.ReadMessage()
		if err != nil {
			return nil, err
		}

		switch b[0] {
		case ClientHbType:
			// logger.Println("heartbeat")
			continue
		case LogoutType:
			// logger.Println("logout")
			return nil, nil
		default:
			// logger.Println("unknown type")
			return nil, fmt.Errorf("unknown format %s", string(b[0]))
		}
	}

}

func (s *serverConn) Close() error {
	close(s.resetC)
	return s.rwc.Close()
}

func (s *serverConn) Reader() io.Reader {
	r, ok := s.sr.(io.Reader)
	if !ok {
		return nil
	}
	return r
}

// options
type ServerOptionFunc func(srv *serverConn)

func SetAuthandler(fn Authenticator) ServerOptionFunc {
	return func(srv *serverConn) {
		srv.authFunc = fn
	}
}
