FROM golang:1.25.4-trixie

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build main.go

ENV APP_PORT=3000

EXPOSE 3000

CMD ["./main"]