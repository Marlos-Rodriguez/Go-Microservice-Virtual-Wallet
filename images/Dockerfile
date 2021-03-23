FROM golang:alpine as builder

ENV GO111MODULE=on

WORKDIR /apps/service

COPY go.mod .

COPY go.sum .

RUN go mod download

COPY . .

ENV USER_PORT 3000

RUN go build cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /apps/service .

EXPOSE 3000

CMD ["./main"]