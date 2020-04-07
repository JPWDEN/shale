FROM golang:latest

RUN apt-get update
RUN apt-get install vim -y

ADD . /go/src/github.com/shale
RUN go install github.com/shale

EXPOSE 8080
ENTRYPOINT /go/bin/shale