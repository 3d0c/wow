package server

import (
	"context"
	"fmt"
	"net"

	"github.com/3d0c/wow/pkg/config"
)

type RequestHandlerFunc func(conn net.Conn)

type Server struct {
	la             *net.TCPAddr
	listener       *net.TCPListener
	requestHandler RequestHandlerFunc
	ctx            context.Context
}

func NewServer(ctx context.Context, cfg config.Config) (*Server, error) {
	var (
		la  *net.TCPAddr
		lc  net.ListenConfig
		l   net.Listener
		err error
	)

	if la, err = net.ResolveTCPAddr("tcp", cfg.Addr); err != nil {
		return nil, fmt.Errorf("errorr resolving address '%s' - %s", cfg.Addr, err)
	}

	s := &Server{
		la:  la,
		ctx: ctx,
	}

	if l, err = lc.Listen(s.ctx, "tcp", la.String()); err != nil {
		return nil, fmt.Errorf("error creating listener - %s", err)
	}

	if tcpl, ok := l.(*net.TCPListener); ok {
		s.listener = tcpl
	} else {
		return nil, fmt.Errorf("error asserting to TCPListener")
	}

	return s, nil
}

func (s *Server) SetRequestHandler(fn RequestHandlerFunc) {
	s.requestHandler = fn
}

func (s *Server) Serve() error {
	var (
		conn *net.TCPConn
		err  error
	)

	if s.requestHandler == nil {
		return fmt.Errorf("request function is not set")
	}

	for {
		if conn, err = s.listener.AcceptTCP(); err != nil {
			s.listener.Close()
			return err
		}
		go s.requestHandler(conn)
	}
}
