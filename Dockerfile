#################
# Builder image
#################
FROM alpine:3.16 as base

USER root

RUN apk --no-cache add ca-certificates \
  && update-ca-certificates

RUN addgroup -g 10001 istio-upgrade-consumer && \
    adduser --disabled-password --system --gecos "" --home "/home/istio-upgrade-consumer" --shell "/sbin/nologin" --uid 10001 istio-upgrade-consumer && \
    mkdir -p "/home/istio-upgrade-consumer" && \
    chown istio-upgrade-consumer:0 /home/istio-upgrade-consumer && \
    chmod g=u /home/istio-upgrade-consumer && \
    chmod g=u /etc/passwd

ENV USER=istio-upgrade-consumer
USER 10001
WORKDIR /home/istio-upgrade-consumer

#################
# Builder image
#################
FROM golang:1.22-bullseye AS builder

WORKDIR /app
COPY . .
RUN make build.binaries

#################
# Final image
#################
FROM base

COPY --from=builder /app/bin/istio-upgrade-consumer /usr/local/bin
