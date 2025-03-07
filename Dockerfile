# Build the simulator binary
FROM golang:1.23.2-alpine AS hive_alpine_dependency

RUN apk --no-cache add gcc musl-dev linux-headers build-base

FROM hive_alpine_dependency:latest AS hive_go_dependency

WORKDIR /source
ADD . /source

RUN go mod tidy
RUN cd /source/hiveproxy && go mod tidy
RUN cd /source/simulators/taiko && go mod tidy
RUN rm -rf /source