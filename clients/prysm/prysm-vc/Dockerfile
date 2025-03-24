FROM hive_prysmctl:latest AS prysmctl

FROM gcr.io/prysmaticlabs/prysm/validator:v5.3.1 AS validator

FROM ubuntu:24.10 AS builder
COPY --from=validator /app/cmd/validator/validator /usr/local/bin/validator

WORKDIR /prysm

COPY prysm_vc.sh .
COPY common.sh .

COPY --from=prysmctl /prysmctl/config.yml .

RUN chmod +x prysm_vc.sh
# Generate the version.txt file.
RUN validator --version >/version.txt

ENTRYPOINT ["/prysm/prysm_vc.sh"]
