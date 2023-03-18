package proto

//TODO : Refactor packet implementation into soupbin package

type Reader interface {
	ReadMessage() ([]byte, error)
}

type Writer interface {
	WriteMessage(msg []byte) error
}

type ReaderWriter interface {
	Reader
	Writer
}
