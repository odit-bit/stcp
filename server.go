package stcp

import (
	"errors"
	"log"
	"net"
	"time"
)

type AuthHandlerFunc func(cs *ClienState) LoginError
type DataPublisherFunc func() []byte

type Server struct {
	addr     string
	session  *string
	sequence *int
	in       chan []byte
	out      chan []byte
	errChan  chan error

	listener      net.Listener
	reader        ProtoReader
	writer        ProtoWriter
	dataPublisher DataPublisherFunc
	authHandler   AuthHandlerFunc
}

func newServer() *Server {
	s := &Server{
		addr:          ":6666",
		session:       nil,
		sequence:      nil,
		in:            make(chan []byte, 1),
		out:           make(chan []byte, 1),
		errChan:       make(chan error, 1),
		reader:        nil,
		writer:        nil,
		dataPublisher: nil,
		authHandler:   nil,
	}
	return s
}

func NewServer(addr string, opts ...ServerOptionFunc) *Server {
	s := newServer()
	if len(opts) > 0 {
		for _, fn := range opts {
			fn(s)
		}
	}
	return s
}

func NewServerWithOptions(opt *serverOptions) *Server {
	s := newServer()
	for _, fn := range opt.options {
		fn(s)
	}
	return s
}

func (s *Server) check() error {
	if s.dataPublisher == nil {
		return errors.New("no data publisher func ")
	}
	if s.session == nil {
		return errors.New("no session ")
	}
	return nil
}

func (s *Server) listen(addr string) error {
	//check the necessary evil
	if err := s.check(); err != nil {
		return err
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	s.listener = listener
	return nil
}

func (s *Server) ListenAndServe(addr string) error {
	if err := s.listen(addr); err != nil {
		return err
	}

	// authenticate incoming connection
	conn, _ := s.listener.Accept()
	s.reader = NewReader(conn)
	s.writer = NewWriter(conn)

	return s.handleEvent()

}

func (s *Server) handleEvent() error {
	hb := createHeartBeat(_HEARTBEAT_SERVER_TYPE)
	hbTicker := time.NewTicker(1000 * time.Millisecond)

	go s.readWorker()
	go s.publishWorker()

	for {
		select {
		case msg := <-s.out:
			s.writer.WriteMessage(msg)

		case <-hbTicker.C:
			s.writer.WriteMessage(hb)

		case <-s.in:
			continue

		case err := <-s.errChan:
			return err
		}
	}
}

func (s *Server) readWorker() {
	for {
		b, err := s.reader.ReadMessage()
		if err != nil {
			log.Println(err)
			s.errChan <- err
		}

		t := b[0]
		switch t {
		case _HEARTBEAT_CLIENT_TYPE:
			continue
		case _LOGIN_REQUEST_TYPE:
			if msg, err := s.handleLoginRequest(b); err != nil {
				s.errChan <- err
			} else {
				s.out <- msg
			}
		case _LOGOUT_REQUEST_TYPE:

		default:
			s.errChan <- ErrUnknownPacket
		}
	}
}

func (s *Server) publishWorker() {
	for {
		s.out <- s.dataPublisher()
	}
}

func (s *Server) handleLoginRequest(msg []byte) ([]byte, error) {

	username, password, session, seqnum, err := parseLoginRequestPacket(msg)
	if err != nil {
		return nil, err
	}
	if s.authHandler == nil {
		log.Println("[DEBUG] auth handler not present, authentication skip", username, password, session, seqnum)
	}

	//create response
	response := createLoginAcceptMessage(*s.session, *s.sequence)
	return response, nil
}

// TODO create Config
type ServerOptionFunc func(*Server)

func (s *Server) Opts(opts ...ServerOptionFunc) {
	for _, fn := range opts {
		fn(s)
	}
}

///////////////////

type serverOptions struct {
	options []ServerOptionFunc
}

func NewServerOpts() *serverOptions {
	return &serverOptions{}
}

func (so *serverOptions) Add(fn ServerOptionFunc) {
	so.options = append(so.options, fn)
}

func (so *serverOptions) Sequence(i int) ServerOptionFunc {
	return func(s *Server) {
		s.sequence = &i
	}
}

func (so *serverOptions) Session(str string) ServerOptionFunc {
	return func(s *Server) {
		s.session = &str
	}
}

func (so *serverOptions) DataPublisherFunc(fn DataPublisherFunc) ServerOptionFunc {

	return func(s *Server) {
		s.dataPublisher = fn
	}
}
