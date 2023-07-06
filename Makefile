.PHONY: fmt test bench

fmt:
	go fmt ./...

test:
	go test -cover ./...

bench:
	go test -v -run=. -bench=. -benchtime=1s ./_examples/performance_test.go
	sleep 1s
	go test -v -run=. -bench=. -benchtime=1s ./_examples/pack_test.go

all: fmt test bench