FROM golang:latest

ENV GO111MODULE=on

WORKDIR /apps/service

COPY go.mod .

COPY go.sum .

RUN go mod download

COPY . .

RUN go build cmd/main.go

EXPOSE 3001

CMD ["./main"]