test:
	go test -v -cover ./...
bench:
	go test ./_examples/performance_test.go -v -run=^$$ -bench=. -benchtime=1s
benchrps:
	go test ./_examples/performance_test.go -v -run=^TestServerRPS$$