FROM golang:1.10-stretch

WORKDIR /go/src/app
COPY . .

# Also get test dependencies....
RUN go get -d -v -t ./...
RUN go install -v ./...

# Expose service app ports
EXPOSE 80 81 82

RUN echo $GOROOT
RUN echo $GOPATH

DEBIAN_FRONTEND=noninteractive apt-get -f install -y sqlite3

CMD ["snacktrackerapi"]
