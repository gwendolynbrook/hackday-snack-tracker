# Base python 2.7 container.
FROM python:2.7-stretch
MAINTAINER Gwen Johnson <gwen@maidbot.com>

RUN pip install --upgrade pip
RUN pip install -U evdev requests

# Make a directory for the script and put it there...
RUN mkdir -p /usr/local/scanner_listener
WORKDIR /usr/local/scanner_listener

COPY ./scanner_listener .
RUN groupadd --gid 106 input

# CMD ["sleep", "infinity"]
CMD ["./ui_listen.py"]
# CMD ["python", "-u", "./ui_listen.py"]
