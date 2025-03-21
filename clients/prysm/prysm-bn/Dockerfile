FROM gcr.io/prysmaticlabs/prysm/cmd/prysmctl:v5.3.1 AS prysmctl

FROM gcr.io/prysmaticlabs/prysm/beacon-chain:v5.3.1 AS beacon

FROM ubuntu:24.10
COPY --from=prysmctl /app/cmd/prysmctl/prysmctl /usr/local/bin
COPY --from=beacon /app/cmd/beacon-chain/beacon-chain /usr/local/bin/beacon-chain

WORKDIR /prysm

COPY prysm_bn.sh .
COPY common.sh .
COPY genesis.json .
COPY config.yml .

RUN chmod +x prysm_bn.sh
# Generate the version.txt file.
RUN beacon-chain --version > /version.txt
# Create jwt.hex file
RUN echo "c49690b5a9bc72c7b451b48c5fee2b542e66559d840a133d090769abc56e39e7" > /tmp/jwt.hex

# Export the usual networking ports to allow outside access to the node
EXPOSE 4000 3500 12000/udp 13000

ENTRYPOINT ["/prysm/prysm_bn.sh"]
