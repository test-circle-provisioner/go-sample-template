FROM golang:alpine
ENV SRC github.com/{{ .Org }}/{{ .Name }}
ARG VERSION
COPY . /go/src/${SRC}

RUN apk --update add git gcc \
  && go install \
  -ldflags="-X main.version=$VERSION" \
  ${SRC}/cmd/service \
  && apk del git gcc

ENTRYPOINT ["{{ .Name }}"]
