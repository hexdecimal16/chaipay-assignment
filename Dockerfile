FROM golang:1.18-alpine AS base
WORKDIR /app

ENV GO111MODULE="on"
ENV GOOS="linux"
ENV CGO_ENABLED=0

# RUN apk update \
#   && apk add --no-cache \
#   ca-certificates \
#   curl \
#   tzdata \
#   git \
#   && update-ca-certificates


FROM base AS dev
WORKDIR /app

RUN go install github.com/cosmtrek/air@latest && go install github.com/go-delve/delve/cmd/dlv@latest
EXPOSE 5000
EXPOSE 2345

ENTRYPOINT ["air"]

FROM base AS builder
WORKDIR /app

COPY . /app
RUN go mod download \
  && go mod verify

RUN go build -o chaipay-assignment -a .

FROM alpine:latest as prod

COPY --from=builder /app/chaipay-assignment /usr/local/bin/chaipay-assignment
EXPOSE 5000

ENTRYPOINT ["/usr/local/bin/chaipay-assignment"]