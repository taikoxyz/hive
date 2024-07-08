ARG baseimage=us-docker.pkg.dev/evmchain/images/taiko-client
ARG tag=taiko-client-v0.29.0

FROM $baseimage:$tag

# Generate the version.txt file.
RUN taiko-client --version > /version.txt

WORKDIR /taiko-client

ADD jwt.hex .
ADD driver.sh .
ADD common.sh .

ENTRYPOINT ["sh", "/taiko-client/driver.sh"]