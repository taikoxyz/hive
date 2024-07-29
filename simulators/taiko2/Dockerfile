# Build the simulator binary
FROM golang:1.22-alpine AS builder
RUN apk --no-cache add gcc musl-dev linux-headers cmake make clang build-base clang-static clang-dev

WORKDIR /source
ADD . /source
RUN go mod download && go build -o /usr/local/bin/sim


# Build within simulator folder
WORKDIR /source
RUN go build -o ./ .

# Build the runner container.
FROM alpine:latest
COPY --from=builder /usr/local/bin/sim /usr/local/bin/sim

ENTRYPOINT ["sim"]
