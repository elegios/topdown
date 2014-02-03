package main

import (
	"github.com/elegios/topdown/server/core"
	"net/http"
)

const (
	rootdir       = "."
	clientdir     = "client"
	fspath        = "/"
	websocketpath = "/ws"
	host          = ":9000"
)

func main() {
	core.Load(rootdir)
	core.Start(fspath, websocketpath, clientdir)
	http.ListenAndServe(host, nil)
}
