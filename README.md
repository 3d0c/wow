## “Word of Wisdom” tcp server.

Design and implement “Word of Wisdom” tcp server.

- TCP server should be protected from DDOS attacks with the Proof of Work (https://en.wikipedia.org/wiki/Proof_of_work), the challenge-response protocol should be used.
- The choice of the POW algorithm should be explained.
- After Proof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other collection of the quotes.
- Docker file should be provided both for the server and for the client that solves the POW challenge

### Getting started

There is a simple Makefile with following targets:

```sh
# Start server
make start-server

# Start client
make start-client

# Build server image
make build-server

# Build client image
make build-client
```

### Binary protocol

Client Server communications are using simple binary protocol:

First 8 bytes are two uint32. First 4 bytes is a length of payload, second 4 bytes is a type of message. Then a payload is following.

```sh
 Length   Type        Payload
00000001 00000011 010011110101010101001010.....
```

### PoW

A Hashcash PoW has been chosen because of it's simplicity for server to verify a solution. Practically with zero overhead. Also this PoW is well documented and has a simple implementation.
