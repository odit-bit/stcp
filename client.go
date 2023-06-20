package stcp

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"
)

// encapsulate client side implementation
type Client struct {
	//session  string // session number or identifier
	//sequence int    //receiverd sequence packet
	rwc  *Conn
	conf *ClientConfig

	// msgHandler func([]byte) error
}

// A SoupBinTCP connection begins with the client opening a TCP/IP socket to the
// server and sending a Login Request Packet
func Connect(conf *ClientConfig) (*Client, error) {

	/*NEW IMPLEMENT*/
	// tcp, err := net.Dial("tcp", conf.Address)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	cli := &Client{
		conf: conf,
	}

	// login
	if err := cli.login(); err != nil {
		return nil, err
	}

	return cli, nil

}

// implement logical packet flow determined soupbinTCP

func (cli *Client) login() error {
	conn, err := net.Dial("tcp", cli.conf.Address)
	if err != nil {
		log.Fatal(err)
	}
	cli.rwc = NewConn(conn)

	err = sendLoginRequest(cli.rwc,
		cli.conf.Username,
		cli.conf.Password,
		cli.conf.session,
		cli.conf.sequenceStr())

	if err != nil {
		return err
	}

	sess, seq, err := receivedLoginResponse(cli.rwc)
	if err != nil {
		return err
	}

	cli.conf.session = sess
	cli.conf.sequence = int64(seq)

	go cli.rwc.StartHeartbeats(ClientHbType)
	return nil
}

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
		//loginAccept
		session = string(bytes.TrimSpace(msg[1:7]))
		sequenceNum, err = strconv.Atoi(string(bytes.TrimSpace(msg[7:17])))

	case LoginRejectType:
		//loginReject
		err = fmt.Errorf("login rejected with reason code %v", string(msg[1]))
	}
	return session, sequenceNum, err

}

func (cli *Client) readMessage() ([]byte, error) {

	for {
		msg, err := cli.rwc.ReadMessage()
		if err != nil {
			log.Println("try to relog", err)
			err = cli.rwc.Close()
			if err != nil {
				log.Println("failed terminate current connection", err)
				return nil, err
			}

			time.Sleep(5 * time.Second)
			err = cli.login()
			if err != nil {
				log.Println("relogin failed", err)
				return nil, err
			}

			log.Println("relogin success", err)
			continue
		}

		switch msg[0] {
		case ServerHbType:
			fmt.Println("hb")
			continue
		case SequenceType:
			return msg, nil
		case SessionEnd:
			return nil, ErrReceiveEndSession
		default:
			return nil, ErrUnknownPacket
		}
	}
}

// return sequence message type
func (cli *Client) ReadSequenceMessage() ([]byte, error) {
	msg, err := cli.readMessage()
	if err != nil {
		return nil, err
	}
	cli.conf.sequence++
	return msg, err
}

func (cli *Client) Close() error {
	return cli.rwc.Close()
}

type ClientConfig struct {
	Address  string
	Username string
	Password string
	session  string
	sequence int64
}

func (conf *ClientConfig) sequenceStr() string {
	num := strconv.Itoa(int(conf.sequence))
	return num
}
