ARG baseimage=us-docker.pkg.dev/evmchain/images/taiko-client
ARG tag=latest

FROM $baseimage:$tag

# Generate the version.txt file.
RUN taiko-client --version >/version.txt

WORKDIR /taiko-client

ADD common.sh .
ADD prover.sh .
ADD jwt.hex /tmp/jwt.hex

ENTRYPOINT ["sh", "/taiko-client/prover.sh"]