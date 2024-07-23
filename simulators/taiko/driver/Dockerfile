FROM golang:1.21-alpine as builder

WORKDIR /taiko
COPY go.mod go.sum *.go ./
RUN go mod download && go build -o /usr/local/bin/taiko


FROM alpine:latest

WORKDIR /taiko

COPY .env .
COPY l1contract_txs.txt .
COPY jwtsecret .

COPY --from=builder /usr/local/bin/taiko /usr/local/bin/taiko

ENTRYPOINT ["taiko"]