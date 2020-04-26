FROM alpine:3.7
# MAINTAINER JunYing Gao <gaojunying@yjy.com>

RUN apk update && \
    apk upgrade && \
    apk --no-cache add curl jq file libc6-compat

VOLUME [ /root ]
WORKDIR /root
EXPOSE 26656 26657
ENTRYPOINT ["/usr/bin/wrapper.sh"]
CMD ["start"]
STOPSIGNAL SIGTERM

COPY wrapper.sh /usr/bin/wrapper.sh