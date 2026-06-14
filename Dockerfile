FROM alpine:3.21

ARG TARGETPLATFORM

RUN apk add --no-cache ca-certificates tzdata \
 && adduser -D -H -u 10001 pixiv

COPY $TARGETPLATFORM/pixiv /usr/bin/pixiv

USER pixiv

ENTRYPOINT ["/usr/bin/pixiv"]
