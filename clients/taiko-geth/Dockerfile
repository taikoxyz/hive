ARG baseimage=us-docker.pkg.dev/evmchain/images/taiko-geth
ARG tag=taiko

FROM $baseimage:$tag as builder

FROM alpine:latest
COPY --from=builder /usr/local/bin/geth /usr/local/bin/geth

# Inject the startup script.
ADD geth.sh /geth.sh
ADD jwt.hex /jwt.hex
RUN chmod +x /geth.sh

# Export the usual networking ports to allow outside access to the node
EXPOSE 8545 8546 8551 30303 30303/udp

ENTRYPOINT ["/geth.sh"]
