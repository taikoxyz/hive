ARG baseimage=us-docker.pkg.dev/evmchain/images/taiko-reth
ARG tag=v1.0.0

FROM $baseimage:$tag AS builder

# Generate the version.txt file.
RUN taiko-reth --version | head -1 | head -1 > /version.txt

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
