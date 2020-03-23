FROM golang:alpine AS builder

WORKDIR /app

RUN apk update --no-cache && \
    apk add git && \
    rm -rf /var/cache/apk/*

COPY go.mod /app
COPY go.sum /app

# Get dependancies - will be cached if we don't change mod/sum
RUN go mod download

ADD ./ /app

ENV GO111MODULE=on

# cgo disabled to allow binary to be self-contained (i.e dependencies)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o golang-test .

# Second stage empty base image
FROM scratch

EXPOSE 8000

# Copy binary from builder stage
WORKDIR /app
COPY --from=builder /app/golang-test .

ENTRYPOINT ["/app/golang-test"]
