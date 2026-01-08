.PHONY: fmt test bench

all: fmt test

fmt:
	go fmt ./...

test:
	go test -v -cover ./...

bench:
	go test -v -run=. -bench=. -benchtime=1s ./_examples/performance_test.go

benchpack:
	go test -v -run=. -bench=. -benchtime=1s ./_examples/pack_test.go
