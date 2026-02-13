.PHONY: fmt test bench

all: fmt test

fmt:
	go fmt ./...

test:
	go test -v -cover ./...

bench:
	go test -v -run=. -bench=. -benchmem -benchtime=1s ./_examples/packet_test.go
