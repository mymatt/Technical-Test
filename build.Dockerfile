
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
