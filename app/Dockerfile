FROM golang:alpine AS builder

WORKDIR /app

RUN apk update --no-cache && \
    apk add git && \
    rm -rf /var/cache/apk/*

ADD ./ /app

ENV GO111MODULE=on

RUN go build -o golang-test .

ENTRYPOINT ["/app/golang-test"]

EXPOSE 8000
