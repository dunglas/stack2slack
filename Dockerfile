FROM golang:alpine

RUN mkdir -p /go/src/github.com/dunglas/stack2slack
WORKDIR /go/src/github.com/dunglas/stack2slack
ADD . /go/src/github.com/dunglas/stack2slack

RUN apk add --no-cache --virtual .build-deps git \
    && go get \
    && go install \
    && apk del .build-deps \
    && rm -rf /go/pkg \
    && rm -rf /go/src \
    && rm -rf /go/cache/apk/*

ENTRYPOINT ["/go/bin/stack2slack"]
