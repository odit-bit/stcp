package stcp

import "errors"

var ErrUnknownPacket = errors.New("unknown packet")
var ErrNotAuthorize = errors.New("not authorize")
var ErrSessionUnavailabe = errors.New("no available session or sequence")
var ErrInvalidPacketField = errors.New("invalid packet field")
var ErrReceiveEndSession = errors.New("receive end session")
var ErrIncompletePacket = errors.New("packet write incomplete")
