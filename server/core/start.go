package core

import (
	"code.google.com/p/go.net/websocket"
	"github.com/elegios/topdown/server/helpers"
	"github.com/elegios/topdown/server/types"
	"log"
	"net/http"
	"os"
)

var (
	slog = log.New(os.Stdout, "     ", log.LstdFlags)
	wlog = log.New(os.Stdout, "WARN ", log.LstdFlags)
)

const (
	defaultActionCount = 1
)

var (
	otm = helpers.NewOneToMany()

	world            *types.World
	defaultCharacter types.Character
)

func Load(root string) {
	world = types.LoadWorld(root)
}

func Start(fspath, wspath, clientdir string) {
	defaultCharacter = types.Character{
		Mapname:   "testmap",
		X:         5,
		Y:         5,
		Id:        "",
		Variation: 0,
		Name:      "",
		Actions:   1,
		Weapon:    "",
		Armor:     "",
		Health:    10,
		MaxHealth: 10,
		Inventory: nil,
		ViewDist:  11,
	}

	//Things that will stay
	ch := make(chan request, 1)
	go otm.Listen()
	go center(ch)

	http.Handle(wspath, websocket.Handler(clientHandler(ch)))
	http.Handle(fspath, http.FileServer(http.Dir(clientdir)))

	slog.Printf("Core has started. clientDir: %#v, wspath: %#v, fspath: %#v\n", clientdir, wspath, fspath)
}
