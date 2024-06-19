FROM us-docker.pkg.dev/evmchain/images/prysm/validator:sha-d54b4dc AS origin

FROM debian:latest

COPY config.yml /
COPY start.sh /

COPY --from=origin /usr/local/bin/validator /usr/local/bin

ENTRYPOINT ["/start.sh"]
