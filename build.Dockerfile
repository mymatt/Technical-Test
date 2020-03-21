
FROM golang:alpine AS builder

# arguments assign env
ARG vers
ENV VERS=$vers

ARG desc
ENV DESC=$desc

ARG sha
ENV SHA=$sha

ENV GOPATH=/go

ENV PATH="$PATH:/go/bin"

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

# Import user and group from the builder stage
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Use non-privileged user
USER usergo:usergo
