ARG baseimage=postgres
ARG tag=16

FROM $baseimage:$tag AS builder

# Generate the version.txt file.
RUN postgres --version > /version.txt

ENV POSTGRES_USER=blobscan
ENV POSTGRES_DB=blobscan_dev
ENV POSTGRES_PASSWORD=s3cr3t