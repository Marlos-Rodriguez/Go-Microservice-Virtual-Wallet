FROM golang:latest

ENV GO111MODULE=on

WORKDIR /apps/service

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

ENV USER_PORT 3002

RUN go build cmd/main.go

EXPOSE 3002

CMD ["./main"]