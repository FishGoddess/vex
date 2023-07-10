.PHONY: fmt test bench

fmt:
	go fmt ./...

test:
	go test -count=1 -cover ./...

bench:
	go test -v -run=. -bench=. -benchtime=1s ./_examples/performance_test.go

benchpack:
	go test -v -run=. -bench=. -benchtime=1s ./_examples/pack_test.go

all: fmt test bench