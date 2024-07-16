ARG baseimage=us-docker.pkg.dev/evmchain/images/taiko-client
ARG tag=taiko-client-v0.29.0

FROM $baseimage:$tag

# Generate the version.txt file.
RUN taiko-client --version >/version.txt

WORKDIR /taiko-client

ADD common.sh .
ADD prover.sh .
ADD jwt.hex .

ENTRYPOINT ["sh", "/taiko-client/prover.sh"]