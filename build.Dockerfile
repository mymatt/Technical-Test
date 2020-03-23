
FROM golang:alpine AS builder

ENV GOPATH=/go

# Setup New User env variables
ENV USER=usergo
ENV UID=10001

# Create new non-privileged User
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nohome" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

WORKDIR /go/

COPY . .

WORKDIR /go/goapi

# Get dependancies - will be cached if we don't change mod/sum
RUN go mod download

# install github.com/gorilla/mux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install -v

FROM scratch

# arguments assign env
ARG vers
ENV VERS=$vers

ARG desc
ENV DESC=$desc

ARG sha
ENV SHA=$sha

ENV PATH="$PATH:/go/bin"

EXPOSE 8080

# Import user and group from the builder stage
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Use non-privileged user
USER usergo:usergo

# Copy binary from builder stage
WORKDIR /go/bin
COPY --from=builder /go/bin .

ENTRYPOINT ["goapi"]
