package main

import "github.com/segmentio/serve/logger"
import "github.com/segmentio/go-log"
import "github.com/tj/docopt"
import "net/http"

var Version = "1.1.0"

const Usage = `
  Usage:
    serve <dir> [--bind addr] [--prefix path]
    serve -h | --help
    serve --version

  Options:
    -p, --prefix path   url prefix [default: /]
    -b, --bind addr     bind address [default: 0.0.0.0:3000]
    -h, --help          output help information
    -v, --version       output version

`

func main() {
	args, err := docopt.Parse(Usage, nil, true, Version, false)
	log.Check(err)

	prefix := args["--prefix"].(string)
	addr := args["--bind"].(string)
	dir := args["<dir>"].(string)

	log.Info("binding to %s", addr)
	log.Info("serving %s", dir)

	h := logger.New(http.StripPrefix(prefix, http.FileServer(http.Dir(dir))))

	err = http.ListenAndServe(addr, h)
	if err != nil {
		log.Error("failed to bind: %s", err)
	}
}
