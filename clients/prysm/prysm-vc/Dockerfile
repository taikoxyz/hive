ARG baseimage=gcr.io/prysmaticlabs/prysm/validator
ARG tag=v5.3.1

FROM $baseimage:$tag as validator

FROM ubuntu:24.10
COPY --from=validator /app/cmd/validator/validator /usr/local/bin/validator

WORKDIR /prysm

COPY prysm_vc.sh .
COPY common.sh .
COPY config.yml .

RUN chmod +x prysm_vc.sh
# Generate the version.txt file.
RUN validator --version >/version.txt

ENTRYPOINT ["/prysm/prysm_vc.sh"]
