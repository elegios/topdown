package main

import (
	"github.com/elegios/topdown/server/core"
	"net/http"
)

const (
	worlddir      = "world"
	clientdir     = "client"
	fspath        = "/"
	websocketpath = "/ws"
	host          = ":9000"
)

func main() {
	core.Load(worlddir)
	core.Start(fspath, websocketpath, clientdir)
	go core.CliControl()
	http.ListenAndServe(host, nil)
}
