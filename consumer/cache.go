package consumer

// represent in-memory data store for message received from stcp-connection
type MessageCache struct {
	Session  string
	Sequence uint64
	Data     map[uint64][]byte
}

func NewCache(session string) *MessageCache {
	mc := MessageCache{
		Session:  session,
		Sequence: 1,
		Data:     map[uint64][]byte{},
	}
	return &mc

}

func (mc *MessageCache) Store(key uint64, value []byte) error {
	mc.Data[key] = value
	mc.Sequence = key
	return nil
}
