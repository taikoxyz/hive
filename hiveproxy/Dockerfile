FROM hive_go_dependency:latest AS builder

# Get dependencies first. This improves caching behavior since they only need
# to be re-downloaded when go.mod changes.
WORKDIR /source
ADD . /source

# Now build the proxy executable.
RUN go build -o /bin/hiveproxy ./tool

# Pull the executable into a fresh image.
FROM hive_alpine_dependency:latest
COPY --from=builder /bin/hiveproxy .
EXPOSE 8081/tcp
ENTRYPOINT ./hiveproxy --addr :8081
