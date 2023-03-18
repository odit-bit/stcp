package stcp

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"strconv"
)

const (
	_LOGIN_REQUEST_TYPE  = 'L'
	_LOGIN_ACCEPT_TYPE   = 'A'
	_LOGIN_REJECT_TYPE   = 'J'
	_LOGOUT_REQUEST_TYPE = 'O'
)

const (
	_LOGIN_REQUEST_LEN  = 1 + 6 + 10 + 10 + 20
	_LOGIN_ACCEPT_LEN   = 1 + 10 + 20
	_LOGIN_REJECT_LEN   = 1 + 1
	_LOGOUT_REQUEST_LEN = 1
)

// login reject reason not authorize
const REJECT_AUTH = 'A'

// login reject reason session unavailable
const REJECT_SESSION = 'S'

type Authenticator interface {
	Auth(username, password string, cs *ClienState) LoginError
}

type LoginError error
type ErrodCodeN uint8

var rejectCode map[ErrodCodeN]LoginError = map[ErrodCodeN]LoginError{
	REJECT_AUTH:    ErrNotAuthorize,
	REJECT_SESSION: ErrSessionUnavailabe,
}

func LoginErrorCode(code ErrodCodeN) LoginError {
	r := rejectCode[code]
	return r
}

type ClienState struct {
	Username string
	Password string
	Session  string
	Sequence int
}

func clientAuth(msg []byte, auth Authenticator) (*ClienState, error) {
	username, password, session, seqnum, err := parseLoginRequestPacket(msg)
	if err != nil {
		return nil, errors.Join(errors.New("failed parse login request"), err)
	}

	cs := ClienState{
		Session:  session,
		Sequence: seqnum,
	}
	err = auth.Auth(username, password, &cs)
	if err != nil {
		return nil, err

	}

	return &cs, nil

}

func loginResponse(msg []byte) (*ClienState, error) {
	switch msg[0] {
	case _LOGIN_ACCEPT_TYPE:
		cs, err := parseLoginAccept(msg)
		if err != nil {
			return nil, err
		}
		log.Println("login accept -", cs)
		return cs, err
	case _LOGIN_REJECT_TYPE:
		return nil, parseLoginRejectMessage(msg)
	default:
		return nil, ErrUnknownPacket
	}

}

// ///////////////////////////////////////
func parseLoginAccept(msg []byte) (*ClienState, error) {
	var session string
	var sequence int

	session = string(bytes.TrimSpace(msg[1:11]))
	seq := string(bytes.TrimSpace(msg[11:31]))
	if seq == "" {
		sequence = 0
	}
	sequence, err := strconv.Atoi(seq)
	if err != nil {
		return nil, ErrInvalidPacketField
	}

	cs := ClienState{
		Session:  session,
		Sequence: sequence,
	}
	return &cs, nil
}

func parseLoginRequestPacket(msg []byte) (username, password, session string, seqnum int, err error) {
	username = string(bytes.TrimSpace(msg[1:7]))
	password = string(bytes.TrimSpace(msg[7:17]))
	session = string(bytes.TrimSpace(msg[17:27]))
	seqnum, err = strconv.Atoi(string(bytes.TrimSpace(msg[27:47])))
	if err != nil {
		err = fmt.Errorf(fmt.Sprintf("%v, got %v", err.Error(), msg[27:47]))
	}
	return username, password, session, seqnum, err
}

func parseLoginRejectMessage(msg []byte) LoginError {
	code := msg[1]
	return LoginErrorCode(ErrodCodeN(code))
}

//////////////////////////////////////////////////////////////

func createLoginRejectMessage(code uint8) []byte {
	b := make([]byte, 2)
	b[0] = _LOGIN_REJECT_TYPE
	b[1] = code
	return b
}

func createLoginAcceptMessage(session string, sequence int) []byte {

	msg := make([]byte, 31)
	msg[0] = _LOGIN_ACCEPT_TYPE

	space := []byte{' '}
	leftPadding(session, space, msg[1:11])
	leftPadding(strconv.Itoa(sequence), space, msg[11:31])

	return msg
}

func createLoginRequestMessage(username, password, session string, sequence int) []byte {
	msg := make([]byte, _LOGIN_REQUEST_LEN)
	msg[0] = _LOGIN_REQUEST_TYPE

	space := []byte{' '}
	leftPadding(username, space, msg[1:7])
	leftPadding(password, space, msg[7:17])
	leftPadding(session, space, msg[17:27])
	leftPadding(strconv.Itoa(sequence), space, msg[27:47])

	return msg
}

// ////////////////////////////////////////////////////////////////

/*DEPRECATE*/
// func leftPadSpace(s string, size int) []byte {
// 	space := []byte{32}
// 	b := bytes.Repeat(space, size)
// 	copy(b, []byte(s))
// 	return b
// }

func leftPadding(s string, pad []byte, buf []byte) {
	offset := len(buf) - len(s)
	copy(buf[:len(s)], []byte(s))

	padding := bytes.Repeat(pad, offset)
	copy(buf[len(s):], padding)

}
