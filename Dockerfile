FROM golang:1.5.2

RUN go get gopkg.in/yaml.v2
RUN go get github.com/jteeuwen/go-pkg-rss

COPY . /go/src/github.com/nicr9/hathor
WORKDIR /go/src/github.com/nicr9/hathor

RUN go install -v

CMD hathor
