FROM golang:1.25.4-trixie

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download