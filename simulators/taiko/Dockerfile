# Build the simulator binary
FROM golang:1.23.2-alpine AS hive_alpine_dependency

RUN apk --no-cache add gcc musl-dev linux-headers build-base

FROM hive_alpine_dependency:latest AS hive_taiko_dependency

WORKDIR /source
ADD . /source

RUN go mod tidy && rm -rf /source

FROM hive_taiko_dependency:latest AS hive_taiko_binary

WORKDIR /source
ADD . /source

RUN go mod tidy && go build -v

# Build the runner container.
FROM hive_alpine_dependency:latest
COPY --from=hive_taiko_binary /source/taiko /bin/taiko
COPY params/jwt.hex /tmp/

ENTRYPOINT ["taiko"]
