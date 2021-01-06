FROM golang:latest

ENV GO111MODULE=on

WORKDIR /apps/service

COPY go.mod .

COPY go.sum .

RUN go mod download

COPY . .

ENV USER_PORT 3000

RUN go build cmd/main.go

EXPOSE 3000

CMD ["./main"]