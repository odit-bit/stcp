package stcp


// type Prefix struct {
// 	Length uint16
// 	Typ    uint8
// }

// type LoginRequest struct {
// 	Prefix
// 	Username [6]byte
// 	Password [10]byte
// 	Session  [10]byte
// 	Sequence [20]byte
// }

// func (lr *LoginRequest) ReadFrom(src io.Reader) (int64, error) {
// 	if err := binary.Read(src, binary.BigEndian, lr); err != nil {
// 		return 0, err
// 	}

// 	return int64(49), nil
// }

// func (lr *LoginRequest) WriteTo(dst io.Writer) (int64, error) {
// 	if err := binary.Write(dst, binary.BigEndian, lr); err != nil {
// 		return 0, err
// 	}
// 	return int64(49), nil
// }

// type LoginAccept struct {
// 	Prefix
// 	Session   [10]byte
// 	Sequenced [20]byte
// }

// type LoginReject struct {
// 	Prefix
// 	Reason uint8
// }

// type Sequenced struct {
// 	Prefix
// 	Payload []byte
// }

// type Unsequenced struct {
// 	Prefix
// 	Payload []byte
// }

//
//convinient constructor function concerned with packet_type
//

// func NewLoginRequest(username, password, session, sequence string) *LoginRequest {
// 	lr := LoginRequest{
// 		Prefix: Prefix{
// 			Length: 1 + 46,
// 			Typ:    'L',
// 		},
// 		Username: [6]byte([]byte(fmt.Sprintf("%-6s", username))),
// 		Password: [10]byte([]byte(fmt.Sprintf("%-10s", password))),
// 		Session:  [10]byte([]byte(fmt.Sprintf("%-10s", session))),
// 		Sequence: [20]byte([]byte(fmt.Sprintf("%-20s", sequence))),
// 	}
// 	return &lr
// }

// // return slice of string consist request data
// // [0]username, [1]password, [2]session [3]sequence
// func ParseLoginRequest(lr *LoginRequest) []string {
// 	login := []string{}

// 	username := bytes.TrimSpace(lr.Username[:])
// 	login = append(login, string(username))

// 	password := bytes.TrimSpace(lr.Password[:])
// 	login = append(login, string(password))

// 	session := bytes.TrimSpace(lr.Session[:])
// 	login = append(login, string(session))

// 	sequence := bytes.TrimSpace(lr.Sequence[:])
// 	login = append(login, string(sequence))

// 	return login
// }

// func NewLoginReject(reason uint8) *LoginReject {
// 	lj := LoginReject{
// 		Prefix: Prefix{
// 			Length: 1 + 1,
// 			Typ:    LoginRejectType,
// 		},
// 		Reason: reason,
// 	}
// 	return &lj
// }

// func ParseLoginReject(lj *LoginReject) string {
// 	return string(lj.Reason)
// }

// func NewUnsequenced(data []byte) *Unsequenced {
// 	us := Unsequenced{
// 		Prefix: Prefix{
// 			Length: uint16(1 + len(data)),
// 			Typ:    'U',
// 		},
// 		Payload: data,
// 	}
// 	return &us
// }

// func NewSequenced(data []byte) *Sequenced {
// 	us := Sequenced{
// 		Prefix: Prefix{
// 			Length: uint16(1 + len(data)),
// 			Typ:    'S',
// 		},
// 		Payload: data,
// 	}
// 	return &us
// }

// func NewLoginAccept(sessionNum string, sequenceNum string) *LoginAccept {
// 	// loginAccepted
// 	la := &LoginAccept{
// 		Prefix: Prefix{
// 			Length: uint16(1 + 30),
// 			Typ:    LoginAcceptType,
// 		},
// 		Session:   [10]byte([]byte(fmt.Sprintf("%-10s", sessionNum))),
// 		Sequenced: [20]byte([]byte(fmt.Sprintf("%-20s", sequenceNum))),
// 	}
// 	return la
// }

// func ParseLoginAccept(la *LoginAccept) []string {
// 	session := string(bytes.TrimSpace(la.Session[:]))
// 	Sequenced := string(bytes.TrimSpace(la.Sequenced[:]))

// 	return append([]string{}, session, Sequenced)
// }
