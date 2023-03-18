package soupbin

//this type is mandatory to make conversation with soupbinTCP-implemented protocol

import (
	"bytes"
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

type loginSession struct {
	session  []byte
	sequence []byte
}

func (ls *loginSession) Session() string {
	return string(ls.session)
}

func (ls *loginSession) SequenceNumber() int {
	seq, err := strconv.Atoi(string(ls.sequence))
	if err != nil {
		log.Printf("[DEBUG] %v \n", err)
	}
	return seq
}

//////////////////////////////////////////////////////////

func parseLoginAccept(msg []byte) (*loginSession, error) {

	ls := loginSession{
		session:  bytes.TrimSpace(msg[1:11]),
		sequence: bytes.TrimSpace(msg[11:31]),
	}

	return &ls, nil
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

func parseLoginRejectMessage(msg []byte) error {
	return LoginError(msg[1])
}

//////////////////////////////////////////////////////////////

func createLoginRejectMessage(code uint8) []byte {
	b := make([]byte, 2)
	b[0] = _LOGIN_REJECT_TYPE
	b[1] = code
	return b
}

func CreateLoginAcceptMessage(session string, sequence int) []byte {

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

// ////////////////////////////////////////////////////

func leftPadding(s string, pad []byte, buf []byte) {
	offset := len(buf) - len(s)
	copy(buf[:len(s)], []byte(s))

	padding := bytes.Repeat(pad, offset)
	copy(buf[len(s):], padding)

}
