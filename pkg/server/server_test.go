package server

import (
	"context"
	"io"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/3d0c/wow/pkg/config"
)

func TestServer(t *testing.T) {
	var (
		cfg = config.Config{ListenAddr: "127.0.0.1:5050"}
		s   *Server
		err error
	)

	s, err = NewServer(context.Background(), cfg)
	assert.Nil(t, err)
	assert.NotNil(t, s)

	s.SetRequestHandler(echoHandler)

	err = s.Serve()

	assert.Nil(t, err)

}

func echoHandler(conn net.Conn) {
	conn.SetReadDeadline(time.Now().Add(time.Second * 5))
	io.Copy(conn, conn)
	conn.Close()
}
