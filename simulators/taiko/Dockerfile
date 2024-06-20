FROM --platform=linux/amd64 ghcr.io/foundry-rs/foundry:nightly AS foundry

FROM golang:1.21-alpine AS builder

# Build Delve, aims to debug.
RUN go install github.com/go-delve/delve/cmd/dlv@latest

WORKDIR /taiko

COPY . .

# Install dependencies.
RUN apk update && apk add --no-cache gcc #musl-dev linux-headers git make build-base

# Use -gcflags="all=-N -l" aims to debug.
RUN go mod download && go build -gcflags="all=-N -l" -v -o /usr/local/bin/taiko ./

FROM alpine:latest

COPY taiko-mono /taiko-mono
COPY --from=builder /usr/local/bin/taiko /usr/local/bin
COPY --from=foundry /usr/local/bin/forge /usr/local/bin

RUN apk update && apk add --no-cache libstdc++

CMD ["/dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "/usr/local/bin/taiko"]
