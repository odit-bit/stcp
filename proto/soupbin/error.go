package soupbin

import "errors"

var ErrUnknownPacket = errors.New("unknown packet")
var ErrNotAuthorize = errors.New("not authorize")
var ErrSessionUnavailabe = errors.New("no available session or sequence")
var ErrInvalidPacketField = errors.New("invalid packet field")
var ErrReceiveEndSession = errors.New("receive end session")
var ErrIncompletePacket = errors.New("packet write incomplete")

var rejectCode map[uint8]error = map[uint8]error{
	REJECT_AUTH:    ErrNotAuthorize,
	REJECT_SESSION: ErrSessionUnavailabe,
}

//return error according to code
func LoginError(code uint8) error {
	r, ok := rejectCode[code]
	if !ok {
		return errors.New("wrong code")
	}
	return r
}
