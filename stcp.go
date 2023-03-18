package stcp

import (
	"errors"
	"log"
	"net"
	"time"
)

//TODO: HIDE heartbeat implementation

// func runHearbeat(pw ProtoWriter, hb []byte) {
// 	tick := time.NewTicker(1000 * time.Millisecond)

// 	log.Println("Start heartbeat")
// 	for {
// 		select {
// 		case <-tick.C:
// 			err := pw.WriteMessage(hb)
// 			if err != nil {
// 				tick.Stop()
// 				return
// 			}

// 		}
// 	}
// }

// type loginReqPacket struct{}
// type loginAcceptPacekt struct{}

// func handleMessage(p *packet) error {
// 	//while read do not write ,vice versa

// 	switch b[0] {
// 	case _HEARTBEAT_CLIENT_TYPE:
// 		fmt.Println(">>> receive _HEARTBEAT_CLIENT_TYPE ")
// 	case _HEARTBEAT_SERVER_TYPE:
// 		fmt.Println(">>> receive _HEARTBEAT_SERVER_TYPE ")
// 	case _LOGIN_REQUEST_TYPE:
// 		fmt.Println(">>> receive _LOGIN_REQUEST_TYPE ")
// 	case _LOGIN_ACCEPT_TYPE:
// 		fmt.Println(">>> receive _LOGIN_ACCEPT_TYPE ")
// 	case _LOGIN_REJECT_TYPE:
// 		fmt.Println(">>> receive _LOGIN_REJECT_TYPE ")
// 	case _SEQUENCE_TYPE:
// 		fmt.Println(">>> receive _SEQUENCE_TYPE ")
// 	case _LOGOUT_REQUEST_TYPE:
// 		fmt.Println(">>> receive _LOGOUT_REQUEST_TYPE")
// 	case _SESSION_END_TYPE:
// 		fmt.Println(">>> receive _SESSION_END_TYPE")
// 	default:
// 		return ErrUnknownPacket
// 	}

// 	return nil
// }

func readWorker(reader ProtoReader, out chan<- []byte) error {
	for {
		b, err := reader.ReadMessage()
		if err != nil {
			// log.Println("read message error -", err)
			return err
		}
		t := b[0]
		switch t {
		case _HEARTBEAT_SERVER_TYPE:
			continue
		case _SEQUENCE_TYPE:
			out <- b
		default:
			return ErrUnknownPacket
		}
	}
}

type SeqDataHandlerFunc func([]byte) error

func HandleMessage(conn net.Conn, fn SeqDataHandlerFunc) error {
	//1
	if fn == nil {
		return errors.New("no handler for received message")
	}
	//2
	reader := NewReader(conn)
	out := make(chan []byte, 1)
	go readWorker(reader, out)

	//3
	writer := NewWriter(conn)
	hb := CreateHeartBeat('R')
	timer := time.NewTicker(1000 * time.Millisecond)
	// done := false
	for {
		select {
		case msg := <-out:
			err := fn(msg)
			if err != nil {
				return err
			}
		case <-timer.C:
			log.Println("idle")
			err := writer.WriteMessage(hb)
			if err != nil {
				log.Println("write data error :", err)
				return err
			}
		}
	}

}
