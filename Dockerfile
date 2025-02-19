FROM golang:1.22 AS builder

WORKDIR /app

COPY /config/config.yaml /app/config.yaml

COPY go.mod go.sum ./

RUN  go mod download

COPY . .

WORKDIR /app/cmd

RUN  go build -o Payment .

EXPOSE 8081

CMD ["./Payment"]