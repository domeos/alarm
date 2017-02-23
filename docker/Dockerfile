FROM private-registry.sohucs.com/domeos-pub/alpine:3.4

MAINTAINER domeos

# Prepare

COPY glibc-2.23-r3.apk /glibc-2.23-r3.apk

RUN apk add --no-cache bash \
    && apk --no-cache add ca-certificates \
    && apk --no-cache add supervisor \
    && apk add --update openssl \
    && wget -q -O /etc/apk/keys/sgerrand.rsa.pub https://raw.githubusercontent.com/sgerrand/alpine-pkg-glibc/master/sgerrand.rsa.pub \
    && apk add glibc-2.23-r3.apk

RUN mkdir -p /opt/falcon-alarm/logs && cd /opt/falcon-alarm && mkdir /etc/supervisor.d

# Set work-dir

WORKDIR /opt/falcon-alarm

# Copy falcon-alarm files

COPY falcon-alarm ./falcon-alarm
COPY control ./control
COPY cfg.template.json ./cfg.template.json
COPY conf ./conf
COPY static ./static
COPY views ./views

# Copy alarm ini

COPY alarm-supervisor.ini /etc/supervisor.d/alarm-supervisor.ini

# Copy dockerize

COPY dockerize /usr/local/bin/dockerize

# Start service

CMD dockerize -template ./cfg.template.json:./cfg.json supervisord -n
