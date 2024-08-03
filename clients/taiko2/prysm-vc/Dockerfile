ARG baseimage=gcr.io/prysmaticlabs/prysm/validator
ARG tag=stable

FROM $baseimage:$tag as upstream

FROM debian:buster-slim
COPY --from=upstream /app/cmd/validator/validator /usr/local/bin/validator

WORKDIR /prysm

COPY prysm_vc.sh .
COPY common.sh .

RUN validator --version >/version.txt

RUN chmod +x /prysm/prysm_vc.sh

ENTRYPOINT ["/prysm/prysm_vc.sh"]
