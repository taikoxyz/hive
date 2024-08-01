ARG baseimage=gcr.io/prysmaticlabs/prysm/beacon-chain
ARG tag=stable

FROM $baseimage:$tag AS upstream

FROM debian:buster-slim
COPY --from=upstream /app/cmd/beacon-chain/beacon-chain /usr/local/bin/beacon-chain

WORKDIR /prysm

COPY prysm_bn.sh .

RUN beacon-chain --version > /version.txt

# Export the usual networking ports to allow outside access to the node
EXPOSE 4000

ENTRYPOINT ["bash", "/prysm/prysm_bn.sh"]
