FROM docker.io/library/alpine:3.16 as runtime

RUN \
  apk add --update --no-cache \
    bash \
    curl \
    ca-certificates \
    tzdata

ENTRYPOINT ["provider-exoscale"]
CMD ["operator"]
COPY provider-exoscale /usr/bin/

USER 65536:0
