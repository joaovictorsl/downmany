FROM golang:1.23.1 as setup

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN mkdir dataset
RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/downmany

CMD ["sleep infinity"]
