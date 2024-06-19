ARG baseimage=us-docker.pkg.dev/evmchain/images/taiko-client
ARG tag=taiko-client-v0.26.0

FROM $baseimage:$tag as builder

FROM alpine:latest
COPY --from=builder /usr/local/bin/taiko-client /usr/local/bin/taiko-client

RUN apk update && apk add --no-cache \
    libstdc++

# Generate the version.txt file.
RUN /usr/local/bin/taiko-client --version > /version.txt

ADD common.sh taiko-client.sh driver.sh proposer.sh prover.sh jwt.hex /
RUN chmod +x /taiko-client.sh

ENTRYPOINT ["sh", "/taiko-client.sh"]
