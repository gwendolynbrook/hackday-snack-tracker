FROM golang:1.10-stretch

RUN apt-get -f install -y curl
RUN apt-get update && apt-get -f install -y sqlite3

WORKDIR /go/src/app
COPY ./snack_tracker_api .

# Also get test dependencies....
RUN go get -d -v -t ./...
RUN go install -v ./...

# Expose service app ports
EXPOSE 80 81 82

RUN echo $GOROOT
RUN echo $GOPATH

CMD ["snacktrackerapi"]
