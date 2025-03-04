ARG baseimage=blossomlabs/blobscan-indexer
ARG tag=latest

FROM --platform=linux/amd64 $baseimage:$tag AS builder

RUN echo "blobscan-indexer latest" >/version.txt
