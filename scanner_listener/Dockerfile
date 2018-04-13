# Base python 2.7 container.
FROM python:2.7-stretch
MAINTAINER Gwen Johnson <gwen@maidbot.com>

RUN apt-get update && apt-get -f install -y python-dev python-pip gcc
# RUN apt-get -f install -y linux-headers-$(uname -r)

RUN pip install -qU --ignore-installed evdev requests

# Make a directory for the script and put it there...
RUN adduser --disabled-password --gecos '' scanner_listener_user
RUN mkdir -p /usr/local/scanner_listener
WORKDIR /usr/local/scanner_listener

COPY ./scanner_listener .

CMD ["./ui_listen.py"]
