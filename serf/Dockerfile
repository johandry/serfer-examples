FROM alpine:3.6

RUN apk add --update wget zip unzip \
    && rm /var/cache/apk/* \
    && echo

# Install JQ to view pretty JSON files
# ADD https://github.com/stedolan/jq/releases/download/jq-1.5/jq-linux64 /usr/local/bin/jq
# RUN chmod +x /usr/local/bin/jq

ADD https://releases.hashicorp.com/serf/0.8.1/serf_0.8.1_linux_amd64.zip /serf.zip
RUN unzip serf.zip \
    && rm serf.zip \
    && chmod +x /serf \
    && mv /serf /usr/local/bin/serf

EXPOSE 7946 7373

COPY entrypoint.sh /
RUN chmod +x /entrypoint.sh

ENTRYPOINT [ "/entrypoint.sh" ]