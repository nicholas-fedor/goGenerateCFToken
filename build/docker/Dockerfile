ARG BUILDPLATFORM
ARG TARGETPLATFORM
ARG BASE_IMAGE=alpine:3.24.0@sha256:660e0827bd401543d81323d4886abbd08fda0fe3ba84337837d0b11a67251283

# Builder Stage
FROM --platform=$BUILDPLATFORM ${BASE_IMAGE} AS builder
RUN apk add --no-cache ca-certificates tzdata

# Final Build
FROM scratch

ARG TARGETPLATFORM

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

COPY --chmod=0755 ${TARGETPLATFORM}/gogeneratecftoken /gogeneratecftoken

USER 65532:65532

ENTRYPOINT ["/gogeneratecftoken"]
