package stcp

import (
	"errors"
	"fmt"
)

var ErrUnknownPacket = errors.New("unknown packet")
var ErrNotAuthorize = errors.New("not authorize")
var ErrSessionUnavailabe = errors.New("no available session or sequence")
var ErrInvalidPacketField = errors.New("invalid packet field")
var ErrReceiveEndSession = errors.New("receive end session")
var ErrIncompletePacket = errors.New("packet write incomplete")

var ErrReadBuffer = fmt.Errorf("not sufficient buffer length for read ")

type OpError struct {
	arg   string
	value string
}

func newOpErr(arg, value string) OpError {
	ce := OpError{
		arg:   arg,
		value: value,
	}
	return ce
}

func (ce *OpError) Error() string {
	errStr := fmt.Sprintf("%v: %v \n", ce.arg, ce.value)
	return errStr
}
