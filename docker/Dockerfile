FROM centos:7

MAINTAINER domeos

# Prepare

RUN rm -f /etc/localtime && ln -s /usr/share/zoneinfo/PRC /etc/localtime && mkdir -p /opt/falcon-alarm && cd /opt/falcon-alarm

# Set work-dir

WORKDIR /opt/falcon-alarm

# Copy falcon-alarm files

COPY falcon-alarm ./falcon-alarm
COPY control ./control
COPY cfg.template.json ./cfg.template.json
COPY conf ./conf
COPY static ./static
COPY views ./views

# Copy start script

COPY run.sh ./run.sh

# Copy dockerize

COPY dockerize /usr/local/bin/dockerize

# Start service

CMD dockerize -template ./cfg.template.json:./cfg.json sh run.sh
