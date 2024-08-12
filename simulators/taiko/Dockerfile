# Build the simulator binary
FROM golang:1.22-alpine AS builder
RUN apk --no-cache add gcc musl-dev linux-headers cmake make clang build-base clang-static clang-dev

WORKDIR /source
ADD . /source

RUN go build -v

# Build within simulator folder
WORKDIR /source
RUN go build -o ./ .

# Build the runner container.
#FROM golang:1.22-alpine
#COPY --from=builder /source/taiko /usr/local/bin/taiko

ENTRYPOINT ["/source/taiko"]
