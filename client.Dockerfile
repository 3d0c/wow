FROM golang:1.22.3 AS builder

WORKDIR /build

COPY . .

RUN go mod tidy

RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/client

# multistage build to copy only binary and config
FROM scratch

COPY --from=builder /build/main /

ENTRYPOINT ["/main"]
