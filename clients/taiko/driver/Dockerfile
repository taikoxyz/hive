ARG baseimage=us-docker.pkg.dev/evmchain/images/taiko-client
ARG tag=p2p-payloadQueue

FROM $baseimage:$tag

# Generate the version.txt file.
RUN taiko-client --version > /version.txt
# Create jwt.hex file
RUN echo "c49690b5a9bc72c7b451b48c5fee2b542e66559d840a133d090769abc56e39e7" > /tmp/jwt.hex

WORKDIR /taiko-client

ADD driver.sh .
ADD common.sh .

# expose soft server port.
EXPOSE 7001/tcp 9222/tcp 9222/udp

ENTRYPOINT ["sh", "/taiko-client/driver.sh"]