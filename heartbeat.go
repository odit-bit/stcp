package stcp

const (
	_HEARTBEAT_SERVER_TYPE = 'H'
	_HEARTBEAT_CLIENT_TYPE = 'R'
)

const (
// _HEARTBEAT_SERVER_LEN = 1
// _HEARTBEAT_CLIENT_LEN = 1
)

func createHeartBeat(typ uint8) []byte {
	b := make([]byte, 1)
	b[0] = typ
	return b
}

func CreateHeartBeat(typ uint8) []byte {
	return createHeartBeat(typ)
}
