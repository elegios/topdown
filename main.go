package main

import (
	"code.google.com/p/go.net/websocket"
	"log"
	"net/http"
	"os"
)

const (
	defaultActionCount = 1
	idlength           = 7
	websocketpath      = "/ws"
	host               = ":9000"
)

var (
	defaultCharacter Character
	otm              = newOneToMany()
	warn             = log.New(os.Stdout, "WARN ", log.LstdFlags)
)

func main() {
	//Temporary loading of things
	defaultCharacter = Character{
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
		viewDist:  11,
	}
	maps["testmap"] = parseMap("testmap.png")

	//Things that will stay
	go otm.Listen()
	ch := make(chan request, 1)
	go Center(ch)
	http.Handle(websocketpath, websocket.Handler(clientHandler(ch)))
	clientDir := "client"
	if len(os.Args) >= 2 {
		clientDir = os.Args[1]
	}
	http.Handle("/", http.FileServer(http.Dir(clientDir)))
	log.Printf("About to serve %s as well as websockets.\n", clientDir)
	http.ListenAndServe(host, nil)
}

var (
	itemData   = make(map[string]ItemData)
	characters = make(map[string]*Character)
	items      = make(map[string][]Item)
	props      = make(map[string][]Prop)
	maps       = make(map[string][][]bits)
)

type Item struct {
	X  int    `json:"x"`
	Y  int    `json:"y"`
	Id string `json:"id"`
}

type Prop struct {
	X         int  `json:"x"`
	Y         int  `json:"y"`
	Variation int  `json:"variation"`
	Collide   bool `json:"collide"`
}

type ItemData struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Variation   string `json:"variation"`
	Description string `json:"description"`
}

type Character struct {
	Mapname   string   `json:"mapname"`
	X         int      `json:"x"`
	Y         int      `json:"y"`
	Id        string   `json:"id"`
	Variation int      `json:"variation"`
	Name      string   `json:"name"`
	Actions   int      `json:"actions"`
	Weapon    string   `json:"weapon"`
	Armor     string   `json:"armor"`
	Health    int      `json:"health"`
	MaxHealth int      `json:"maxhealth"`
	Inventory []string `json:"inventory"`
	viewDist  int
}
