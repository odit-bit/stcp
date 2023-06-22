package stcp

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"
)

// encapsulate client side implementation to login and receive sequence message
type Client struct {
	rwc  *Conn
	conf *ClientConfig

	seqMsgC chan []byte
	errC    chan error
}

// A SoupBinTCP connection begins with the client opening a TCP/IP socket to the
// server and sending a Login Request Packet
func Connect(conf *ClientConfig) (*Client, error) {
	cli := &Client{
		rwc:     &Conn{},
		conf:    conf,
		seqMsgC: make(chan []byte),
		errC:    make(chan error),
	}

	// login
	if err := cli.login(); err != nil {
		return nil, err
	}

	go cli.start()
	return cli, nil
}

// implementation of login mechanism in soupbinTCP
func (cli *Client) login() error {
	conn, err := net.Dial("tcp", cli.conf.Address)
	if err != nil {
		log.Fatal(err)
	}
	cli.rwc = NewConn(conn)

	err = sendLoginRequest(cli.rwc,
		cli.conf.Username,
		cli.conf.Password,
		cli.conf.Session,
		cli.conf.sequenceStr())

	if err != nil {
		return err
	}

	sess, seq, err := receivedLoginResponse(cli.rwc)
	if err != nil {
		return err
	}

	cli.conf.Session = sess
	cli.conf.Sequence = int64(seq)

	return nil
}

func (cli *Client) start() {
	go cli.rwc.StartHeartbeats(ClientHbType)

	//continously read from connection
	for {
		msg, err := cli.rwc.ReadMessage()
		if err != nil {
			log.Println("try to relog", err)
			err = cli.rwc.Close()
			if err != nil {
				log.Println("failed terminate current connection", err)
				// return nil, err
			}

			//todo: refactor this line
			time.Sleep(5 * time.Second)

			err = cli.login()
			if err != nil {
				log.Println("relogin failed", err)
				// return nil, err
			}

			log.Println("relogin success", err)
			continue
		}

		switch msg[0] {
		case ServerHbType:
			//connection still live
			continue

		case SequenceType:
			cli.seqMsgC <- msg
			continue

		case SessionEnd:
			cli.errC <- ErrReceiveEndSession

		default:
			cli.errC <- ErrUnknownPacket
		}
		break
	}

}

// return sequence message byte
func (cli *Client) Receive() ([]byte, error) {

	select {
	case msg := <-cli.seqMsgC:
		cli.conf.Sequence++
		return msg, nil
	case err := <-cli.errC:
		return nil, err
	}

}

func (cli *Client) SequenceNum() int {
	return int(cli.conf.Sequence)
}

func (cli *Client) Close() error {
	close(cli.errC)
	close(cli.seqMsgC)
	return cli.rwc.Close()
}

// convinient function to make message from args and write to w
func sendLoginRequest(w SequenceWriter, username, password, session, sequence string) error {
	lr := []byte{}
	lr = append(lr, []byte(fmt.Sprintf("%-6s", username))...)
	lr = append(lr, []byte(fmt.Sprintf("%-10s", password))...)
	lr = append(lr, []byte(fmt.Sprintf("%-10s", session))...)
	lr = append(lr, []byte(fmt.Sprintf("%-20s", sequence))...)

	if err := w.WriteMessage(LoginRequestType, lr); err != nil {
		return err
	}
	return nil
}

func receivedLoginResponse(sr SequenceReader) (string, int, error) {
	//retreive login response
	var session string
	var sequenceNum int
	var err error
	msg, err := sr.ReadMessage()
	if err != nil {
		return session, sequenceNum, err
	}

	switch msg[0] {
	case LoginAcceptType:
		session = string(bytes.TrimSpace(msg[1:7]))
		sequenceNum, err = strconv.Atoi(string(bytes.TrimSpace(msg[7:17])))

	case LoginRejectType:
		err = fmt.Errorf("login rejected with reason code %v", string(msg[1]))
	}
	return session, sequenceNum, err

}

type ClientConfig struct {
	Address  string
	Username string
	Password string
	Session  string
	Sequence int64
}

func (conf *ClientConfig) sequenceStr() string {
	num := strconv.Itoa(int(conf.Sequence))
	return num
}
