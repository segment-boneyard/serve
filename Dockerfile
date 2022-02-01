FROM 528451384384.dkr.ecr.us-west-2.amazonaws.com/segment-golang:1.17.6

RUN go get github.com/segmentio/go-log
RUN go get github.com/tj/docopt
RUN go get github.com/dustin/go-humanize


ADD . /go/src/github.com/segmentio/serve
RUN go install github.com/segmentio/serve

COPY mime.types /etc/mime.types

ENTRYPOINT ["/go/bin/serve"]
CMD [".", "--bind", "0.0.0.0:3000"]

EXPOSE 3000
