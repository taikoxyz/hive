FROM us-docker.pkg.dev/evmchain/images/prysm/prysmctl:taiko AS prysmctl
FROM us-docker.pkg.dev/evmchain/images/prysm/beacon-chain:taiko AS beacon
FROM us-docker.pkg.dev/evmchain/images/prysm/validator:taiko AS validator
FROM us-docker.pkg.dev/evmchain/images/client-go:devnet AS geth

FROM debian:latest

RUN apt-get update && apt-get install -y \
  procps && \
  rm -rf /var/lib/apt/lists/*


COPY --from=prysmctl /usr/local/bin/prysmctl /usr/local/bin
COPY --from=beacon /usr/local/bin/beacon-chain /usr/local/bin
COPY --from=validator /usr/local/bin/validator /usr/local/bin
COPY --from=geth /usr/local/bin/geth /usr/local/bin

# Generate the version.txt file.
RUN /usr/local/bin/geth console --exec 'console.log(admin.nodeInfo.name)' --maxpeers=0 --nodiscover --dev > logs.out 2>/dev/null | head -1 > /version.txt

WORKDIR /devnet
RUN mkdir -p /devnet/consensus /devnet/execution/keystore


COPY start.sh .
COPY geth.sh .
COPY beacon-chain.sh .
COPY validator.sh .
COPY config.yml         /devnet/consensus
COPY jwtsecret          /devnet/consensus
COPY genesis.json       /devnet/execution
COPY geth_password.txt  /devnet/execution
COPY jwtsecret          /devnet/execution
COPY keyfile.json       /devnet/execution/keystore

# Export the usual networking ports to allow outside access to the node
EXPOSE 4000 8545 8546 8551 30303 30303/udp

ENTRYPOINT ["sh", "/devnet/start.sh"]