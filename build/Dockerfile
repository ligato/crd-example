FROM golang:1.11-alpine as build
ENV PACKAGEPATH=github.com/ligato/crd-example/
RUN apk add --update protobuf git bash gcc musl-dev
COPY [".","/go/src/${PACKAGEPATH}"]
WORKDIR /go/src/${PACKAGEPATH}/
RUN ./scripts/build.sh --race-test-disabled
ENTRYPOINT ["/go/bin/crd"]
