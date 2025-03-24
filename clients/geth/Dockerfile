FROM hive_prysmctl:latest AS prysmctl

FROM ethereum/client-go:v1.15.5 AS geth

FROM ubuntu:24.10
COPY --from=geth /usr/local/bin/geth /usr/local/bin

WORKDIR /geth

COPY common.sh .
COPY geth.sh .
COPY keyfile.json .

COPY --from=prysmctl /prysmctl/genesis.json .

RUN chmod +x geth.sh
# Generate the version.txt file.
RUN geth console --exec 'console.log(admin.nodeInfo.name)' --maxpeers=0 --nodiscover --dev > logs.out 2>/dev/null | head -1 > /version.txt
# Create jwt.hex file
RUN echo "c49690b5a9bc72c7b451b48c5fee2b542e66559d840a133d090769abc56e39e7" > /tmp/jwt.hex

# Export the usual networking ports to allow outside access to the node
EXPOSE 8545 8546 8551 30303 30303/udp

ENTRYPOINT ["/geth/geth.sh"]