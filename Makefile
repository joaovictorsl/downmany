build:
	go build -o ./bin/downmany

run: build
	./bin/downmany

test:
	go test ./...
