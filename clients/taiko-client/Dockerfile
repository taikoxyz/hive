ARG baseimage=gcr.io/evmchain/taiko-client
ARG tag=hekla

FROM $baseimage:$tag as builder

FROM alpine:latest
COPY --from=builder /usr/local/bin/taiko-client /usr/local/bin/taiko-client

# Generate the version.txt file.
RUN /usr/local/bin/taiko-client --version > /version.txt

ADD common.sh taiko-client.sh driver.sh proposer.sh prover.sh jwt.hex /
RUN chmod +x /taiko-client.sh

ENTRYPOINT ["sh", "/taiko-client.sh"]
