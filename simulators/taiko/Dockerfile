FROM golang:1.21-alpine as builder

WORKDIR /taiko
COPY go.mod go.sum *.go ./
RUN go mod download && go build -o /usr/local/bin/taiko


FROM alpine:latest

WORKDIR /taiko

COPY .env .
COPY jwtsecret .
COPY scripts/get_env.sh .
COPY l1contract_txs.txt .

COPY --from=builder /usr/local/bin/taiko /usr/local/bin/taiko

ENTRYPOINT ["taiko"]