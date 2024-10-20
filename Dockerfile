FROM golang:1.23.2 as builder

ENV TODO_PORT=7540
ENV TODO_PASSWORD=1234
ENV TODO_DBFILE=data/scheduler.db
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

FROM ubuntu:latest

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    && apt-get install -y --no-install-recommends \
    gcc libc6-dev \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY . .

RUN go mod download

EXPOSE ${TODO_PORT}

RUN go build -o /todoapp

CMD ["./todoapp"]