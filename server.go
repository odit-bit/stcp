package stcp

import (
	"bytes"
	"fmt"
	"log"
	"net"
)

type Server struct {
	rwc *Conn

	//indicate connection failure by not receive expected packet from client
	readErr chan error
}

// If the login request is valid, the server
// responds with a Login Accepted Packet and begins sending Sequenced Data Packets.
// The connection continues until the TCP/IP socket is broken.
func AuthConnection(tcp net.Conn, lv LoginValidator) (*Server, error) {

	// 1.Auth
	rwc, err := acceptConnection(tcp, lv)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	srv := &Server{
		rwc:     rwc,
		readErr: make(chan error),
	}

	go srv.startHearbeats()
	go srv.receive()
	return srv, nil

}

// server only receive unsequence , debug and hearbeat packet from client
func (srv *Server) receive() {
	for {
		msg, err := srv.rwc.ReadMessage()
		if err != nil {
			srv.readErr <- err
			return
		}
		switch msg[0] {
		case ClientHbType:
			fmt.Println("hb")
			continue
		case LogoutType:
			srv.readErr <- ErrReceiveEndSession
			return
		default:
			srv.readErr <- ErrUnknownPacket
			return
		}
	}
}

// sending sequence data packet
func (srv *Server) Send(payload []byte) error {
	//check error
	select {
	case err, ok := <-srv.readErr:
		if !ok {
			return err
		}
	default:
	}

	return srv.rwc.WriteMessage(SequenceType, payload)
}

func (srv *Server) startHearbeats() {
	srv.rwc.StartHeartbeats(ServerHbType)
}

func (srv *Server) EndSession() error {
	err := srv.rwc.WriteMessage(SessionEnd, []byte{})
	if err != nil {
		return err
	}

	return nil
}

func (srv *Server) Close() error {
	return srv.rwc.Close()
}

type ServerConfig struct {
	Address string
	Session string

	LV       LoginValidator
	Producer func() []byte
}

// If the login request is valid, the server
// responds with a Login Accepted Packet and begins sending Sequenced Data Packets.
// The connection continues until the TCP/IP socket is broken.
func acceptConnection(conn net.Conn, lv LoginValidator) (*Conn, error) {

	//NEW implement
	rwc := NewConn(conn)
	ok, err := authConnection(rwc, rwc, lv)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, err
	}

	return rwc, nil
}

func authConnection(sr SequenceReader, sw SequenceWriter, lv LoginValidator) (bool, error) {
	payload, err := sr.ReadMessage() //readPacket(rw)
	if err != nil {
		return false, err
	}

	data := ParseLoginRequestBytes(payload)

	if !lv.IsValid(data[0], data[1]) {
		err := sendLoginReject(sw, LoginRejectAuth)
		return false, err
	}

	sess, seq, ok := lv.SessionValid(data[2], data[3])
	if !ok {
		err := sendLoginReject(sw, LoginRejectSession)
		return false, err
	}

	return true, sendLoginAccept(sw, sess, seq)
}

func sendLoginReject(sw SequenceWriter, reason uint8) error {
	lr := []byte{reason}
	err := sw.WriteMessage(LoginRejectType, lr)
	if err != nil {
		return err
	}

	return nil
}

func sendLoginAccept(sw SequenceWriter, session, sequence string) error {
	msg := []byte{}
	msg = append(msg, []byte(fmt.Sprintf("%-10s", session))...)
	msg = append(msg, []byte(fmt.Sprintf("%-20s", sequence))...)

	if err := sw.WriteMessage(LoginAcceptType, msg); err != nil {
		return err
	}
	return nil
}

// decode loginRequest type byte as slice of string
func LoginRequestByteToSlice(b []byte) ([]string, error) {
	if b[0] != 'L' {
		return nil, fmt.Errorf("unknown type %v", string(b[0]))
	}
	data := ParseLoginRequestBytes(b)
	return data, nil
}

func ParseLoginRequestBytes(p []byte) []string {
	_ = p[46]
	login := []string{}

	username := bytes.TrimSpace(p[1:7])
	login = append(login, string(username))

	password := bytes.TrimSpace(p[7:17])
	login = append(login, string(password))

	session := bytes.TrimSpace(p[17:27])
	login = append(login, string(session))

	sequence := bytes.TrimSpace(p[27:47])
	login = append(login, string(sequence))

	return login
}
