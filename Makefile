build:
	go build -o ./bin/downmany

run-server: build
	./bin/downmany -server

run-client: build
	./bin/downmany -file_hash 1336902055 -server_addr 127.0.0.1:8000

test:
	go test ./...
