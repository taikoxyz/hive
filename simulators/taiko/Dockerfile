FROM golang:1.21-alpine as builder

#RUN apk update && apk add --no-cache gcc musl-dev linux-headers git make build-base

WORKDIR /taiko

COPY go.mod go.sum *.go ./

RUN go mod download

RUN go build -o /usr/local/bin/taiko

FROM ethereum/solc:0.8.24

COPY --from=builder /usr/local/bin/taiko /usr/local/bin/taiko

ENTRYPOINT ["taiko"]