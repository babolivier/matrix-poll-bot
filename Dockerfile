FROM golang:1.12-alpine AS builder
RUN apk add --no-cache git
WORKDIR /opt
COPY . /opt
RUN go build

FROM alpine
COPY --from=builder /opt/matrix-poll-bot /usr/local/bin/
RUN apk add --no-cache su-exec ca-certificates
VOLUME /data
CMD /usr/local/bin/matrix-poll-bot --config /data/config.yaml
