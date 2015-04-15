package main

import "github.com/segmentio/serve/logger"
import "github.com/segmentio/go-log"
import "github.com/tj/docopt"
import "net/http"

var Version = "1.0.0"

const Usage = `
  Usage:
    serve <dir> [--address a]
    serve -h | --help
    serve --version

  Options:
    -a, --address a   bind address [default: localhost:3000]
    -h, --help        output help information
    -v, --version     output version

`

func main() {
	args, err := docopt.Parse(Usage, nil, true, Version, false)
	log.Check(err)

	addr := args["--address"].(string)
	dir := args["<dir>"].(string)

	log.Info("binding to %s", addr)
	log.Info("serving %s", dir)

	h := logger.New(http.FileServer(http.Dir(dir)))

	err = http.ListenAndServe(addr, h)
	if err != nil {
		log.Error("failed to bind: %s", err)
	}
}
