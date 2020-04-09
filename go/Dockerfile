FROM golang:latest

RUN apt-get update
RUN apt-get install vim -y
RUN go get "github.com/go-sql-driver/mysql"

ADD . /go/src/github.com/shale/go
RUN go install github.com/shale/go
RUN ls -al ./bin

EXPOSE 8080
EXPOSE 443
ENTRYPOINT ["./bin/go"]