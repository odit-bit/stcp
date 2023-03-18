package stcp

const _SESSION_END_TYPE = 'Z'

func CreateEndSessionMessage() []byte {
	msg := []byte{_SESSION_END_TYPE}
	return msg
}
