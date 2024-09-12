ARG baseimage=ghcr.io/foundry-rs/foundry
ARG tag=nightly

FROM $baseimage:$tag AS builder

#FROM alpine:3.20.2
#COPY --from=foundry /usr/local/bin/anvil /usr/local/bin

# Generate the version.txt file.
RUN anvil --version > /version.txt

WORKDIR /anvil

COPY anvil.sh .
COPY common.sh .

EXPOSE 8545

ENTRYPOINT ["sh", "/anvil/anvil.sh"]
