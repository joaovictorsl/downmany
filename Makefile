build:
	go build -o ./bin/downmany

run-server: build
	./bin/downmany -server

run-client: build
	./bin/downmany 

test:
	go test ./...
