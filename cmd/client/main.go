package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/3d0c/wow/pkg/config"
	"github.com/3d0c/wow/pkg/pow"
	"github.com/3d0c/wow/pkg/proto"
)

type clientMain struct {
	cfg config.Config
}

func main() {
	var (
		cm   = &clientMain{}
		msg  *proto.Message
		conn net.Conn
		err  error
	)

	flag.StringVar(&cm.cfg.Addr, "addr", "127.0.0.1:5050", "Server address")
	flag.IntVar(&cm.cfg.HashcashMaxIterations, "max_iter", 1000000, "Hashcash max iterations")
	flag.Parse()

	if conn, err = net.Dial("tcp", cm.cfg.Addr); err != nil {
		log.Fatalf("Error connecting to %s\n", cm.cfg.Addr)
	}
	defer func() {
		log.Printf("Closing connection to %s\n", cm.cfg.Addr)
		conn.Close()
	}()

	for {
		if msg, err = cm.handle(conn); err != nil {
			log.Printf("Error handling connection - %s\n", err)
			break
		}

		fmt.Println(msg.String())

		time.Sleep(time.Second * 2)
	}

	return
}

func (cm *clientMain) handle(conn net.Conn) (*proto.Message, error) {
	var (
		msg      *proto.Message
		hashcash pow.HashcashData
		result   pow.HashcashData
		payload  []byte
		err      error
	)

	// Sending challenge request
	msg = proto.NewMessage(proto.RequestChallenge, []byte{})
	if err = msg.Write(conn); err != nil {
		return nil, fmt.Errorf("error writing challenge request - %s", err)
	}

	// Reading challenge response
	if err = msg.Read(conn); err != nil {
		return nil, fmt.Errorf("error reading challenge response - %s", err)
	}
	if err = json.Unmarshal(msg.Payload, &hashcash); err != nil {
		return nil, fmt.Errorf("error unmarshalling response - %s", err)
	}

	if result, err = hashcash.ComputeHashcash(cm.cfg.HashcashMaxIterations); err != nil {
		return nil, fmt.Errorf("error computing hashcash - %s", err)
	}

	if payload, err = json.Marshal(result); err != nil {
		return nil, fmt.Errorf("error marshalling hashcash - %s", err)
	}

	// Sending resource request
	msg = proto.NewMessage(proto.RequestResource, payload)
	if err = msg.Write(conn); err != nil {
		return nil, fmt.Errorf("error writing message - %s", err)
	}

	// Reading resource result
	if err = msg.Read(conn); err != nil {
		return nil, fmt.Errorf("error reading resource response - %s", err)
	}

	return msg, nil
}
