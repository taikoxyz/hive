FROM --platform=linux/amd64 ghcr.io/foundry-rs/foundry:nightly AS builder

FROM alpine:latest

RUN apk update && apk add --no-cache libstdc++

COPY --from=builder /usr/local/bin/forge /usr/local/bin/