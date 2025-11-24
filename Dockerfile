FROM golang:1.25.4-trixie

WORKDIR /app

COPY . .

RUN go mod download

RUN go build main.go

ENV AppPort=3000

EXPOSE 3000

CMD ["./main"]