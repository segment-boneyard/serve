FROM golang:1.8.3

RUN go get github.com/segmentio/go-log
RUN go get github.com/tj/docopt
RUN go get github.com/dustin/go-humanize


ADD . /go/src/github.com/segmentio/serve
RUN go install github.com/segmentio/serve

COPY mime.types /etc/mime.types

ENTRYPOINT ["/go/bin/serve"]
CMD [".", "--bind", "0.0.0.0:3000"]

EXPOSE 3000
