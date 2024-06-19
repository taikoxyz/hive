FROM gcr.io/prysmaticlabs/prysm/cmd/prysmctl@sha256:18eb5c5e40da9c44ee6a691482fb3ae954c26156bab2ed3690aabbd76fb6243c AS prysmctl

FROM us-docker.pkg.dev/evmchain/images/prysm/beacon-chain:sha-d54b4dc AS beacon

FROM debian:latest

RUN apt-get update && \
    apt-get install bc -y && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

COPY --from=prysmctl /prysmctl /usr/local/bin
COPY --from=beacon /usr/local/bin/beacon-chain /usr/local/bin

COPY genesis.json /
COPY config.yml /
COPY jwtsecret /
COPY start.sh /

ENTRYPOINT ["/start.sh"]
