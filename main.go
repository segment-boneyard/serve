package main

import "github.com/segmentio/serve/logger"
import "github.com/segmentio/go-log"
import "github.com/tj/docopt"
import "net/http"

var Version = "0.0.1"

const Usage = `
  Usage:
    serve <dir>
    serve -h | --help
    serve --version

  Options:
    -h, --help       output help information
    -v, --version    output version

`

func main() {
	args, err := docopt.Parse(Usage, nil, true, Version, false)
	log.Check(err)

	addr := "localhost:3000"
	dir := args["<dir>"].(string)

	log.Info("binding to %s", addr)
	log.Info("serving %s", dir)

	h := logger.New(http.FileServer(http.Dir(dir)))

	err = http.ListenAndServe(addr, h)
	if err != nil {
		log.Error("failed to bind: %s", err)
	}
}
