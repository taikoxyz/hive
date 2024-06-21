FROM gcr.io/prysmaticlabs/prysm/cmd/prysmctl@sha256:18eb5c5e40da9c44ee6a691482fb3ae954c26156bab2ed3690aabbd76fb6243c AS prysmctl
FROM --platform=linux/amd64 ghcr.io/foundry-rs/foundry:nightly AS foundry

FROM us-docker.pkg.dev/evmchain/images/client-go:devnet AS geth

FROM alpine:latest

RUN apk update && apk add --no-cache bc

COPY --from=prysmctl /prysmctl /usr/local/bin
COPY --from=foundry /usr/local/bin/forge /usr/local/bin
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
COPY deploy_l1_contract.sh .

# Export the usual networking ports to allow outside access to the node
EXPOSE 8545 8546 8551 30303 30303/udp

ENTRYPOINT ["sh", "/execution/start.sh"]
