# First stage: build the executable.
FROM golang:1.12.5-alpine AS builder

WORKDIR /usr/local/bin

COPY ./cmd /usr/local/bin

ENTRYPOINT [ "./cmd run-server" ]
