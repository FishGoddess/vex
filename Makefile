.PHONY: fmt test bench

fmt:
	go fmt ./...

test:
	go test -cover ./...

bench:
	go test -v -bench=. -benchtime=1s ./_examples/performance_test.go

all: fmt test bench