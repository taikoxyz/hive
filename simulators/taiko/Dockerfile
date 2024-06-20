FROM --platform=linux/amd64 ghcr.io/foundry-rs/foundry:nightly AS foundry

FROM golang:1.21-alpine AS builder

WORKDIR /taiko

COPY . .

# Install dependencies.
RUN apk update && apk add --no-cache gcc #musl-dev linux-headers git make build-base

RUN go mod download && go build -v -o /usr/local/bin/taiko ./

FROM alpine:latest

COPY taiko-mono /taiko-mono
COPY --from=builder /usr/local/bin/taiko /usr/local/bin
COPY --from=foundry /usr/local/bin/forge /usr/local/bin

RUN apk update && apk add --no-cache libstdc++

