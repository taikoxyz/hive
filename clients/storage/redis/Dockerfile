ARG baseimage=redis
ARG tag=7.4.0

FROM $baseimage:$tag AS builder

# Generate the version.txt file.
RUN redis-server -v /version.txt

ENV REDIS_REPLICATION_MODE=master
ENV REDIS_PASSWORD=s3cr3t
ENV REDIS_USERNAME=blobscan