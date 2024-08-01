ARG baseimage=ethereum/client-go
ARG tag=latest

FROM $baseimage:$tag AS builder

FROM debian:latest

COPY --from=geth /usr/local/bin/geth /usr/local/bin

# Generate the version.txt file.
RUN /usr/local/bin/geth console --exec 'console.log(admin.nodeInfo.name)' --maxpeers=0 --nodiscover --dev > logs.out 2>/dev/null | head -1 > /version.txt

WORKDIR /geth

COPY geth.sh .

# Export the usual networking ports to allow outside access to the node
EXPOSE 8545 8546 8551 30303 30303/udp

ENTRYPOINT ["sh", "/geth/geth.sh"]