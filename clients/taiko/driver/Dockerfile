ARG baseimage=us-docker.pkg.dev/evmchain/images/taiko-client
ARG tag=taiko-client-v0.40.0

FROM $baseimage:$tag

# Generate the version.txt file.
RUN taiko-client --version > /version.txt

WORKDIR /taiko-client

ADD driver.sh .
ADD common.sh .

# expose soft server port.
EXPOSE 7000/tcp

ENTRYPOINT ["sh", "/taiko-client/driver.sh"]