package main

import "github.com/segmentio/serve/logger"
import "github.com/segmentio/go-log"
import "github.com/tj/docopt"
import "net/http"
import "fmt"

var Version = "1.1.0"

const Usage = `
  Usage:
    serve <dir> [--bind addr] [--prefix path] [--error path]
    serve -h | --help
    serve --version

  Options:
    -p, --prefix path   url prefix [default: /]
    -b, --bind addr     bind address [default: 0.0.0.0:3000]
    -e, --error path    path to serve on error [default: ]
    -h, --help          output help information
    -v, --version       output version

`

func main() {
	args, err := docopt.Parse(Usage, nil, true, Version, false)
	log.Check(err)

	prefix := args["--prefix"].(string)
	addr := args["--bind"].(string)
	errPage := args["--error"].(string)
	dir := args["<dir>"].(string)

	log.Info("binding to %s", addr)
	log.Info("serving %s", dir)

	h := logger.New(http.StripPrefix(prefix, http.FileServer(http.Dir(dir))))

	err = http.ListenAndServe(addr, NotFoundHook{h: h, ErrorPage: errPage})
	if err != nil {
		log.Error("failed to bind: %s", err)
	}
}

type hookedResponseWriter struct {
	http.ResponseWriter
	Request   *http.Request
	ignore    bool
	ErrorPage string
}

func (hrw *hookedResponseWriter) WriteHeader(status int) {
	if status == 404 && hrw.ErrorPage != "" {
		hrw.ignore = true
		http.Redirect(hrw.ResponseWriter, hrw.Request, hrw.ErrorPage, http.StatusSeeOther)
	} else {
		hrw.ResponseWriter.WriteHeader(status)
	}
}

func (hrw *hookedResponseWriter) Write(p []byte) (int, error) {
	if hrw.ignore {
		return len(p), nil
	}
	return hrw.ResponseWriter.Write(p)
}

type NotFoundHook struct {
	h         http.Handler
	ErrorPage string
}

func (nfh NotFoundHook) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	nfh.h.ServeHTTP(&hookedResponseWriter{ResponseWriter: w, Request: r, ErrorPage: nfh.ErrorPage}, r)
}
