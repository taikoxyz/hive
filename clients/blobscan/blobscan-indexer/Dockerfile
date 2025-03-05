ARG baseimage=blossomlabs/blobscan-indexer
ARG tag=0.3.1

FROM --platform=linux/amd64 $baseimage:$tag AS builder

RUN echo "blobscan-indexer latest" >/version.txt
