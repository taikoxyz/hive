ARG baseimage=ethereum/client-go
ARG tag=v1.14.7

FROM $baseimage:$tag AS builder

FROM debian:latest

COPY --from=builder /usr/local/bin/geth /usr/local/bin

# Generate the version.txt file.
RUN geth console --exec 'console.log(admin.nodeInfo.name)' --maxpeers=0 --nodiscover --dev > logs.out 2>/dev/null | head -1 > /version.txt

WORKDIR /geth

COPY common.sh .
COPY geth.sh .

RUN chmod +x /geth/geth.sh

# Export the usual networking ports to allow outside access to the node
EXPOSE 8545 8546 8551 30303 30303/udp

ENTRYPOINT ["/geth/geth.sh"]