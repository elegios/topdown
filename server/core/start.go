package core

import (
	"code.google.com/p/go.net/websocket"
	"github.com/elegios/topdown/server/helpers"
	"github.com/elegios/topdown/server/script"
	"github.com/elegios/topdown/server/types"
	"log"
	"net/http"
	"os"
)

var (
	slog = log.New(os.Stdout, "     ", log.LstdFlags)
	wlog = log.New(os.Stderr, "WARN ", log.LstdFlags)
)

var (
	otm = helpers.NewOneToMany()

	world            *types.World
	defaultCharacter types.Character
)

func Load(worlddir string) {
	world = script.LoadWorld(worlddir)

	slog.Println("Got maps:")
	for name := range world.Maps {
		slog.Println(name)
	}
}

func Start(fspath, wspath, clientdir string) {
	defaultCharacter = types.Character{
		Pos: types.Position{
			Mapid: "testmap",
			X:     5,
			Y:     5,
		},
		Id:            "",
		Variation:     0,
		Name:          "",
		Actions:       0,
		RecoverySpeed: 1,
		RecoveryMax:   1,
		Weapon:        "",
		Armor:         "",
		Health:        10,
		MaxHealth:     10,
		Inventory:     nil,
		ViewDist:      11,
	}

	//Things that will stay
	ch := make(chan request, 1)
	go otm.Listen()
	go center(ch)

	http.Handle(wspath, websocket.Handler(clientHandler(ch)))
	http.Handle(fspath, http.FileServer(http.Dir(clientdir)))

	slog.Printf("Core has started. clientDir: %#v, wspath: %#v, fspath: %#v\n", clientdir, wspath, fspath)
}
