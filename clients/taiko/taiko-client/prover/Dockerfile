ARG baseimage=us-docker.pkg.dev/evmchain/images/taiko-client
ARG tag=hekla

FROM $baseimage:$tag as builder

FROM alpine:latest

RUN apk add --no-cache --update libstdc++
# Generate the version.txt file.
RUN /usr/local/bin/taiko --version > /version.txt

WORKDIR /prover

ADD common.sh .
ADD prover.sh .

COPY --from=builder /usr/local/bin/taiko-client /usr/local/bin/taiko-client

ENTRYPOINT ["sh", "/prover/prover.sh"]

