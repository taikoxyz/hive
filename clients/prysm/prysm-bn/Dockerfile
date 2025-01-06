ARG baseimage=gcr.io/prysmaticlabs/prysm/beacon-chain
ARG tag=v5.0.4

FROM $baseimage:$tag AS upstream

FROM debian:buster-slim
COPY --from=upstream /app/cmd/beacon-chain/beacon-chain /usr/local/bin/beacon-chain

WORKDIR /prysm

COPY prysm_bn.sh .
COPY common.sh .
COPY jwt.hex /tmp/jwt.hex

RUN beacon-chain --version > /version.txt

RUN chmod +x /prysm/prysm_bn.sh

# Export the usual networking ports to allow outside access to the node
EXPOSE 4000 3500 12000/udp 13000

ENTRYPOINT ["/prysm/prysm_bn.sh"]
