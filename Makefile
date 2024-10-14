start-server:
	go run cmd/server/main.go

start-client:
	go run cmd/client/main.go

build-server:
	docker build -f server.Dockerfile .

build-client:
	docker build -f client.Dockerfile .
