test:
	go test -v -cover ./...
bench:
	go test ./_examples/performance_test.go -v -run=^$$ -bench=^BenchmarkServer$$ -benchtime=1s
benchrps:
	go test ./_examples/performance_test.go -v -run=^TestRPS$$
	sleep 5s
	go test ./_examples/performance_test.go -v -run=^TestRPSWithPool$$