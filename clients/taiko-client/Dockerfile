ARG baseimage=gcr.io/evmchain/taiko-client
ARG tag=hekla

FROM $baseimage:$tag as builder

FROM alpine:latest
COPY --from=builder /usr/local/bin/taiko-client /usr/local/bin/taiko-client