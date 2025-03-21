ARG baseimage=us-docker.pkg.dev/evmchain/images/taiko-reth
ARG tag=latest

FROM us-docker.pkg.dev/evmchain/images/taiko-geth:taiko AS taiko-geth

FROM $baseimage:$tag AS builder
COPY --from=taiko-geth /usr/local/bin/geth /usr/local/bin/geth

RUN apt-get update && apt-get install -y curl jq && apt-get autoclean

# Generate the version.txt file.
RUN taiko-reth --version | head -1 | head -1 > /version.txt
# Create jwt.hex file
RUN echo "c49690b5a9bc72c7b451b48c5fee2b542e66559d840a133d090769abc56e39e7" > /tmp/jwt.hex

WORKDIR /reth

# Inject the startup script.
ADD taiko-reth.sh .

# Inject the enode id retriever script.
RUN mkdir /hive-bin
ADD enode.sh /hive-bin/enode.sh
RUN chmod +x /hive-bin/enode.sh

# Export the usual networking ports to allow outside access to the node
EXPOSE 8545 8546 8551 30303 30303/udp

ENTRYPOINT ["sh", "/reth/taiko-reth.sh"]
