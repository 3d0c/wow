package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"

	"github.com/3d0c/wow/pkg/cache"
	"github.com/3d0c/wow/pkg/config"
	"github.com/3d0c/wow/pkg/pow"
	"github.com/3d0c/wow/pkg/proto"
	"github.com/3d0c/wow/pkg/quotes"
	"github.com/3d0c/wow/pkg/server"
)

type serverMain struct {
	cfg    config.Config
	cache  cache.Cache
	quotes quotes.Quotes
}

func main() {
	var (
		sm  = &serverMain{}
		s   *server.Server
		err error
	)

	flag.StringVar(&sm.cfg.Addr, "listen", "127.0.0.1:5050", "Listen on")
	flag.IntVar(&sm.cfg.Timeout, "timeout", 5, "Read/Write timeout in seconds")
	// flag.IntVar(&sm.cfg.MaxOpen, "maxopen", 1000, "Maximum opened connections")
	flag.Int64Var(&sm.cfg.HashcashDuration, "hash_duration", 300, "Hash duration")
	flag.IntVar(&sm.cfg.ZeroCount, "zero_count", 3, "Count of leading zeros")
	flag.Parse()

	sm.cache = cache.NewInMemoryCache()
	sm.quotes = quotes.NewInMemoryQuotes()

	ctx := context.Background()

	if s, err = server.NewServer(ctx, sm.cfg); err != nil {
		log.Fatalf("Error creating server - %s\n", err)
	}

	s.SetRequestHandler(sm.requestHandler)

	if err = s.Serve(); err != nil {
		log.Fatalf("Error serving - %s\n", err)
	}
}

func (sm *serverMain) requestHandler(conn net.Conn) {
	var (
		remote = conn.RemoteAddr().String()
		msg    = proto.NewMessage(0, []byte{})
		result = proto.NewMessage(0, []byte{})
		err    error
	)

	log.Printf("Starting handler for %s\n", remote)
	defer func() {
		log.Printf("Stopping handler for %s\n", remote)
	}()

	for {
		if err = conn.SetDeadline(time.Now().Add(time.Second * time.Duration(sm.cfg.Timeout))); err != nil {
			log.Printf("Error setting connetion deadline - %s\n", err)
			break
		}

		if err = msg.Read(conn); err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				log.Printf("Timed out - %s\n", err)
			} else {
				log.Printf("Error reading message - %s\n", err)
			}
			break
		}

		if result, err = sm.process(msg, remote); err != nil {
			log.Printf("Error processing request - %s\n", err)
			break
		}

		if err = result.Write(conn); err != nil {
			log.Printf("Error sending response - %s\n", err)
			break
		}
	}

	conn.Close()
	return
}

func (sm *serverMain) process(msg *proto.Message, remote string) (*proto.Message, error) {
	var (
		hashcash pow.HashcashData
		payload  []byte
		err      error
	)

	switch msg.Type {
	case proto.Quit:
		return nil, fmt.Errorf("client %s closed session", remote)

	case proto.RequestChallenge:
		log.Printf("Client %s requested the challenge\n", remote)

		randValue := rand.Intn(100000)
		if err = sm.cache.Add(randValue, sm.cfg.HashcashDuration); err != nil {
			return nil, fmt.Errorf("error adding value into cache - %s", err)
		}

		hashcash = pow.HashcashData{
			Version:    1,
			ZerosCount: sm.cfg.ZeroCount,
			Date:       time.Now().Unix(),
			Resource:   remote,
			Rand:       base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", randValue))),
			Counter:    0,
		}

		if payload, err = json.Marshal(hashcash); err != nil {
			return nil, fmt.Errorf("error marshaling hashcash - %s", err)
		}

		return &proto.Message{
			Type:    proto.ResponseChallenge,
			Payload: payload,
		}, nil

	case proto.RequestResource:
		log.Printf("Client %s requested resource\n", remote)

		if err = json.Unmarshal([]byte(msg.Payload), &hashcash); err != nil {
			return nil, fmt.Errorf("error unmarshaling hashcash - %s", err)
		}

		// validate hashcash params
		if hashcash.Resource != remote {
			return nil, fmt.Errorf("error validating hashcache resource")
		}

		randValueBytes, err := base64.StdEncoding.DecodeString(hashcash.Rand)
		if err != nil {
			return nil, fmt.Errorf("error decoding random value - %s", err)
		}
		randValue, err := strconv.Atoi(string(randValueBytes))
		if err != nil {
			return nil, fmt.Errorf("error parsing random value - %s", err)
		}

		exists, err := sm.cache.Get(randValue)
		if err != nil {
			return nil, fmt.Errorf("error getting random value from cache - %s", err)
		}
		if !exists {
			return nil, fmt.Errorf("challenge expired or not sent")
		}

		if time.Now().Unix()-hashcash.Date > sm.cfg.HashcashDuration {
			return nil, fmt.Errorf("challenge expired")
		}

		maxIter := hashcash.Counter
		if maxIter == 0 {
			maxIter = 1
		}

		if _, err = hashcash.ComputeHashcash(maxIter); err != nil {
			return nil, fmt.Errorf("error computing hashcash - %s", err)
		}

		log.Printf("Client %s succesfully computed hashcash %s\n", remote, string(msg.Payload))

		sm.cache.Delete(randValue)

		return &proto.Message{
			Type:    proto.ResponseResource,
			Payload: sm.quotes.Get(rand.Intn(sm.quotes.Len())),
		}, nil

	default:
		return nil, fmt.Errorf("error processing request - unkown request type")
	}
}
