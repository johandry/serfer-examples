# Build image
FROM golang:1.9.0-alpine3.6 AS build4linux

ARG PKG_NAME=surfer
ARG PKG_BASE=github.com/johandry
ARG BIN_NAME=surfer

ADD . /go/src/${PKG_BASE}/${PKG_NAME}

RUN  cd /go/src/${PKG_BASE}/${PKG_NAME} \
  && GOOS=linux go build -o /${BIN_NAME} \
  && echo

# To do a manual build uncomment this line and comment out the following lines.
# CMD [ "/bin/bash" ]

# Application image
FROM alpine:3.6 AS application

COPY --from=build4linux /surfer .

# Uncomment these lines to install Serf in the image:
# ADD https://releases.hashicorp.com/serf/0.8.1/serf_0.8.1_linux_amd64.zip /serf.zip
# RUN  apk add --update unzip \
#   && rm /var/cache/apk/* \
#   && unzip serf.zip \
#   && rm serf.zip \
#   && chmod +x /serf \
#   && mv /serf /usr/local/bin/serf

EXPOSE 7946 7373

ENTRYPOINT [ "./surfer" ]
