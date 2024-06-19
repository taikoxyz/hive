FROM gcr.io/prysmaticlabs/prysm/cmd/prysmctl@sha256:18eb5c5e40da9c44ee6a691482fb3ae954c26156bab2ed3690aabbd76fb6243c AS prysmctl

FROM us-docker.pkg.dev/evmchain/images/client-go:devnet AS geth

FROM ubuntu:22.04

RUN apt-get update && \
    apt-get install bc -y

COPY --from=prysmctl /prysmctl /usr/local/bin
COPY --from=geth /usr/local/bin/geth /usr/local/bin

# Generate the version.txt file.
RUN /usr/local/bin/geth console --exec 'console.log(admin.nodeInfo.name)' --maxpeers=0 --nodiscover --dev 2>/dev/null | head -1 > /version.txt

WORKDIR /execution

COPY config.yml .
COPY start.sh .
COPY genesis.json .
COPY keyfile.json .
COPY geth_password.txt .
COPY jwtsecret .

ENTRYPOINT ["sh", "/execution/start.sh"]
