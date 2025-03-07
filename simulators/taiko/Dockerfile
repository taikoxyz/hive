# Build the simulator binary
FROM hive_go_dependency:latest AS hive_taiko_binary

WORKDIR /source
ADD . /source

RUN go mod tidy && go build -v

# Build the runner container.
FROM hive_alpine_dependency:latest
COPY --from=hive_taiko_binary /source/taiko /bin/taiko
COPY params/jwt.hex /tmp/

ENTRYPOINT ["taiko"]
