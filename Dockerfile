FROM golang:1.20.10-alpine3.17 AS builder

RUN apk add --update  && \
    apk add --no-cache alpine-conf tzdata git

ADD ./src /go/src/authenticator
ADD ./src/log /go/log
ADD ./src/config /go/config

RUN cd /go/src/authenticator && \
    go install authenticator

FROM alpine:3.18.4 AS app

COPY --from=builder /go/bin/* /go/bin/authenticator
COPY --from=builder /go/log /go/log
COPY --from=builder /go/config /go/config

ENV CONTAINER_NAME=authenticator
ENV HAS_WRITE_LOG_TO_FILE=false
ENV LOG_LEVEL=5

EXPOSE 8080

WORKDIR "/go"
ENTRYPOINT ["/go/bin/authenticator"]