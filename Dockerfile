FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV GOOS linux

RUN apk update --no-cache && apk add --no-cache git musl-dev util-linux-dev gcc

WORKDIR /build

ADD go.mod .

ADD go.sum .

RUN go mod download

COPY . .

RUN go build -a -ldflags "-linkmode external -extldflags '-static' -s -w" -o /app/videocut ./main.go

FROM alpine

RUN apk update --no-cache && apk add --no-cache ca-certificates ffmpeg

WORKDIR /app

COPY --from=builder /app/videocut /app/videocut
