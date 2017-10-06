package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/segmentio/go-log"
	"github.com/segmentio/serve/logger"
	"github.com/tj/docopt"
)

var Version = "1.2.0"

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

	if errPage != "" {
		errPage = dir + "/" + errPage
	}

	server := FileServer(http.Dir(dir), errPage)

	handler := logger.New(http.StripPrefix(prefix, server))

	err = http.ListenAndServe(addr, handler) // NotFoundHook{h: h, ErrorPage: errPage})
	if err != nil {
		log.Error("failed to bind: %s", err)
	}
}

func FileServer(fs http.FileSystem, errorPage string) http.Handler {
	fsh := http.FileServer(fs)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := fs.Open(path.Clean(r.URL.Path))
		if err != nil && !os.IsPermission(err) {
			fmt.Println("Not found: " + r.URL.Path)
			NotFound(fs, w, r, errorPage)
			return
		}
		fsh.ServeHTTP(w, r)
	})
}

func NotFound(fs http.FileSystem, w http.ResponseWriter, r *http.Request, errorPage string) {
	if _, err := os.Stat(errorPage); err == nil {
		w.WriteHeader(http.StatusNotFound)
	}
	WriteFile(fs, w, r, errorPage)
}

func WriteFile(fs http.FileSystem, w http.ResponseWriter, r *http.Request, errorPage string) {
	bys, err := ioutil.ReadFile(errorPage)
	if err != nil {
		msg, code := toHTTPError(err)
		http.Error(w, msg, code)
		return
	}

	w.Write(bys)
}

func toHTTPError(err error) (msg string, httpStatus int) {
	if os.IsNotExist(err) {
		return "404 page not found", http.StatusNotFound
	}
	if os.IsPermission(err) {
		return "403 Forbidden", http.StatusForbidden
	}
	// Default:
	return "500 Internal Server Error", http.StatusInternalServerError
}
