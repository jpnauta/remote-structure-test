FROM golang:latest as build

# Install pre-requisites
RUN apt-get update && apt-get install make
RUN go get github.com/golang/dep/cmd/dep

WORKDIR /go/src/github.com/jpnauta/remote-structure-test

# Install deps
ADD Gopkg.* ./
RUN dep ensure --vendor-only

# Build app
ADD cmd/ cmd/
ADD pkg/ pkg/
ADD Makefile .
RUN make cross && \
    ln -s /go/src/github.com/jpnauta/remote-structure-test/out/remote-structure-test-linux-amd64 /usr/bin/remote-structure-test

# Prepare for tests
COPY docker/scripts/run-unit-tests.sh /
RUN chmod +x /run-unit-tests.sh
ADD tests/ tests/

FROM alpine:latest as runtime
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=build /usr/bin/remote-structure-test /usr/bin/remote-structure-test

ENTRYPOINT ["remote-structure-test"]
