package stcp

import "time"

const (
	//default deadline to timeout the connection
	TimeoutDefault = 4 * time.Second

	// interval or ticker for send heartbeat packet
	HeartbeatTimeout = 2 * time.Second
)
const (
	SequenceType     uint8 = 'S'
	UnsesequenceType uint8 = 'U'
	DebugType        uint8 = '+'

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
