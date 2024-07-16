ARG baseimage=us-docker.pkg.dev/evmchain/images/taiko-client
ARG tag=sha256:2e4ba363a115708fad4bd0c0899b1bab8b9454dde753ec6e489b3b417d5ca1ed

FROM $baseimage@$tag

# Generate the version.txt file.
RUN taiko-client --version >/version.txt

WORKDIR /taiko-client

ADD common.sh .
ADD prover.sh .
ADD jwt.hex .

ENTRYPOINT ["sh", "/taiko-client/prover.sh"]