package stcp

import (
	"bytes"
	"fmt"
)

// decode loginRequest type byte as slice of string
func LoginRequestByteToSlice(b []byte) ([]string, error) {
	if b[0] != 'L' {
		return nil, fmt.Errorf("unknown type %v", string(b[0]))
	}
	data := ParseLoginRequestBytes(b)
	return data, nil
}

func ParseLoginRequestBytes(p []byte) []string {
	_ = p[46]
	login := []string{}

	username := bytes.TrimSpace(p[1:7])
	login = append(login, string(username))

	password := bytes.TrimSpace(p[7:17])
	login = append(login, string(password))

	session := bytes.TrimSpace(p[17:27])
	login = append(login, string(session))

	sequence := bytes.TrimSpace(p[27:47])
	login = append(login, string(sequence))

	return login
}
