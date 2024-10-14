package proto

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

const (
	Quit              = iota // on quit each side (server or client) should close connection
	RequestChallenge         // from client to server - request new challenge from server
	ResponseChallenge        // from server to client - message with challenge for client
	RequestResource          // from client to server - message with solved challenge
	ResponseResource         // from server to client - message with useful info is solution is correct, or with error if not
)

var (
	BytesOrder = binary.BigEndian
)

type Message struct {
	Type    uint32
	Payload []byte
}

func NewMessage(typ uint32, payload []byte) *Message {
	return &Message{
		Type:    typ,
		Payload: payload,
	}
}

// Read first 8 bytes and treat first 4 as a payload size, second 4 as request type
func (m *Message) Read(conn net.Conn) error {
	var (
		err error
	)

	header := make([]byte, 8)
	if _, err = io.ReadFull(conn, header); err != nil {
		return fmt.Errorf("error reading header - %s", err)
	}

	payloadSz := BytesOrder.Uint32(header[:4])

	m.Type = BytesOrder.Uint32(header[4:])
	m.Payload = make([]byte, payloadSz)

	if _, err = io.ReadFull(conn, m.Payload); err != nil {
		return fmt.Errorf("error reading payload - %s", err)
	}

	return nil
}

func (m *Message) Write(conn net.Conn) error {
	var (
		buf = bytes.NewBuffer(nil)
		sz  = uint32(len(m.Payload))
		err error
	)

	if err = binary.Write(buf, BytesOrder, sz); err != nil {
		return fmt.Errorf("error composing message to send (length) - %s", err)
	}
	if err = binary.Write(buf, BytesOrder, m.Type); err != nil {
		return fmt.Errorf("error composing message to send (type) - %s", err)
	}
	if err = binary.Write(buf, BytesOrder, m.Payload); err != nil {
		return fmt.Errorf("error composing message to send (payload) - %s", err)
	}

	if _, err = conn.Write(buf.Bytes()); err != nil {
		return fmt.Errorf("error writing message - %s", err)
	}

	return nil
}

func (m *Message) Len() int {
	return 4 + len(m.Payload)
}

func (m *Message) String() string {
	return string(m.Payload)
}
