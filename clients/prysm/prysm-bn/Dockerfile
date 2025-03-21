FROM hive_prysmctl:latest AS prysmctl

FROM gcr.io/prysmaticlabs/prysm/beacon-chain:v5.3.1 AS beacon

FROM ubuntu:24.10 AS builder
COPY --from=beacon /app/cmd/beacon-chain/beacon-chain /usr/local/bin/beacon-chain

WORKDIR /prysm

COPY prysm_bn.sh .
COPY common.sh .

COPY --from=prysmctl /prysmctl/config.yml .
COPY --from=prysmctl /prysmctl/genesis.ssz .
COPY --from=prysmctl /prysmctl/genesis.json .

RUN chmod +x /prysm/prysm_bn.sh
# Generate the version.txt file.
RUN beacon-chain --version > /version.txt
# Create jwt.hex file
RUN echo "c49690b5a9bc72c7b451b48c5fee2b542e66559d840a133d090769abc56e39e7" > /tmp/jwt.hex

# Export the usual networking ports to allow outside access to the node
EXPOSE 4000 3500 12000/udp 13000

ENTRYPOINT ["/prysm/prysm_bn.sh"]
