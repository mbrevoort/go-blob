FROM golang:1.4.2

ADD . /go/src/github.com/mbrevoort/go-blob
RUN go install github.com/mbrevoort/go-blob
CMD go-blob

EXPOSE 3000