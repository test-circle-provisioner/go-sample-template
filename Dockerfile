FROM golang:alpine
ENV SRC github.com/segmentio/go-hello-world
ARG VERSION
COPY . /go/src/${SRC}

RUN apk --update add git gcc \
  && go install \
  -ldflags="-X main.version=$VERSION" \
  ${SRC}/cmd/go-hello-world \
  && apk del git gcc

ENTRYPOINT ["go-hello-world"]
