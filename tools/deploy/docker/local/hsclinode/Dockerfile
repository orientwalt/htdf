FROM alpine:3.7
# MAINTAINER JunYing Gao <gaojunying@yjy.com>

RUN apk update && \
    apk upgrade && \
    apk --no-cache add curl jq file

VOLUME [ /root ]
WORKDIR /root

EXPOSE 1317 26656 26657
