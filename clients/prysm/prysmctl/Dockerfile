FROM gcr.io/prysmaticlabs/prysm/cmd/prysmctl:v5.3.1 AS prysmctl

FROM ubuntu:24.10
COPY --from=prysmctl /app/cmd/prysmctl/prysmctl /usr/local/bin/prysmctl

WORKDIR /prysmctl
COPY genesis.json .
COPY config.yml .

RUN prysmctl  \
    testnet  \
    generate-genesis \
    --fork=fulu \
    --num-validators=1 \
    --genesis-time-delay=0 \
    --output-ssz=genesis.ssz \
    --chain-config-file=config.yml \
    --geth-genesis-json-in=genesis.json \
    --geth-genesis-json-out=genesis.json

RUN rm /usr/local/bin/prysmctl