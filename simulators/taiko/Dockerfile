FROM --platform=linux/amd64 ghcr.io/foundry-rs/foundry:nightly AS foundry

FROM golang:1.21-alpine as builder
#RUN apk update && apk add --no-cache gcc musl-dev linux-headers git make build-base

WORKDIR /taiko
COPY go.mod go.sum *.go ./
RUN go mod download && go build -o /usr/local/bin/taiko


FROM ethereum/solc:0.8.24

WORKDIR /taiko

RUN mkdir -p /taiko/taiko-mono/packages/

COPY deploy_l1_contract.sh .
COPY taiko-mono/packages/protocol /taiko/taiko-mono/packages/protocol

COPY --from=builder /usr/local/bin/taiko /usr/local/bin/taiko
COPY --from=foundry /usr/local/bin/forge /usr/local/bin
COPY --from=foundry /usr/local/bin/cast /usr/local/bin

ENTRYPOINT ["taiko"]