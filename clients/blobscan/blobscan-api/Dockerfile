ARG baseimage=blossomlabs/blobscan-api
ARG tag=1.5.0

FROM --platform=linux/amd64 $baseimage:$tag AS builder

RUN echo "blobscan-indexer latest" >/version.txt
