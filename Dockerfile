FROM segment/sources-build:v1.0.0 as builder
WORKDIR /go/src/github.com/segmentio/go-hello-world/
COPY . .
RUN govendor install -ldflags '-s -w' .

FROM debian:stretch
RUN apt-get update -y && apt-get install -y ca-certificates
COPY --from=builder /go/bin/go-hello-world /go-hello-world
ENTRYPOINT [ "/goHelloWorld" ]