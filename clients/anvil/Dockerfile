FROM hive_prysmctl:latest AS prysmctl

FROM ghcr.io/foundry-rs/foundry:nightly

# Generate the version.txt file.
#RUN anvil --version > /version.txt

WORKDIR /anvil

COPY anvil.sh .
COPY common.sh .

COPY --from=prysmctl /prysmctl/genesis.json .

EXPOSE 8545

ENTRYPOINT ["sh", "/anvil/anvil.sh"]
