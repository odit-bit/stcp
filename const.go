package stcp

import "time"

const (
	// WriteTimeoutDefault = 5 * time.Second
	ReadTimeoutDefault = 3 * time.Second
	HeartbeatTimeout   = 2 * time.Second

	DefaultBuffer = 4096
)
const (
	SequenceMessage     uint8 = 'S'
	UnsesequenceMessage uint8 = 'U'
	DebugMessage        uint8 = '+'

	LoginRequestType uint8 = 'L'
	LoginAcceptType  uint8 = 'A'

	LoginRejectType    uint8 = 'R'
	LoginRejectAuth    uint8 = 'A'
	LoginRejectSession uint8 = 'S'

	LogoutType uint8 = 'O'

	SessionEnd uint8 = 'Z'

	ServerHbType uint8 = 'H'
	ClientHbType uint8 = 'R'
)
