package main

import (
	"code.google.com/p/go.net/websocket"
	"log"
	"sync"
)

//this data is shared, message relating to a
//character may be accepted even though the
//character has been killed, as knowledge of
//death will not be propagated until a tick,
//and a message at the same time might even
//be a data-race. It means the server has to
//check that the id really corresponds to a
//character.
type clientData struct {
	charIds map[string]struct{}
	sync.Mutex
}

type mapPos struct {
	X    int  `json:"x"`
	Y    int  `json:"y"`
	Data bits `json:"data"`
}

type tickUpdate struct {
	Maps         map[string][]mapPos `json:"maps"`
	Controllable []string            `json:"controllable"`
	Characters   []*Character        `json:"characters"`
	Props        []Prop              `json:"props"`
	Items        []Item              `json:"items"`
}

type request struct {
	message map[string]string
	ch      chan map[string]string
}

func clientHandler(ch chan<- request) func(*websocket.Conn) {
	return func(conn *websocket.Conn) {
		//TODO: load owned characters
		data := &clientData{
			charIds: make(map[string]struct{}),
		}
		go sendTicks(conn, data)
		locChan := make(chan map[string]string, 1)
		for {
			var mess map[string]string
			err := websocket.JSON.Receive(conn, &mess)
			if err != nil {
				warn.Println("A websocket connection died:", err)
				return
			}

			if !ok(mess, data) {
				warn.Print("A message was not ok:", mess)
				continue
			}

			ch <- request{mess, locChan}
			switch mess["command"] {
			case "create":
				res := <-locChan
				log.Printf("got return from create")
				data.Lock()
				data.charIds[res["id"]] = struct{}{}
				data.Unlock()

			case "itemrequest":
				websocket.JSON.Send(conn, <-locChan)
			}
		}
	}
}

func ok(mess map[string]string, data *clientData) bool {
	if _, ok := mess["character"]; ok {
		data.Lock()
		defer data.Unlock()
		_, ok = data.charIds[mess["character"]]
		return ok
	}
	return true
}

func sendTicks(conn *websocket.Conn, data *clientData) {
	c := make(chan struct{}, 1)
	otm.AddP <- Pair{conn, c}
	defer func() { otm.Rem <- conn }()
	for {
		<-c
		log.Println("Preparing to send tick to", conn.LocalAddr())
		update := tickUpdate{
			Maps:         make(map[string][]mapPos),
			Controllable: make([]string, 0),
			Characters:   make([]*Character, 0),
			Props:        make([]Prop, 0),
			Items:        make([]Item, 0),
		}
		updateControllable(data, &update)
		collectMapPos(&update)
		collectVisible(&update)
		//done collecting data
		otm.Done()
		err := websocket.JSON.Send(conn, &update)
		if err != nil {
			warn.Println("ticksender died:", err)
			return
		}
	}
}

func updateControllable(data *clientData, update *tickUpdate) {
	data.Lock()
	defer data.Unlock()
	for charId := range data.charIds {
		_, ok := characters[charId]
		if !ok {
			delete(data.charIds, charId)
			continue
		}
		update.Controllable = append(update.Controllable, charId)
	}
}

func collectMapPos(update *tickUpdate) {
	for _, cid := range update.Controllable {
		char := characters[cid]
		j := max(char.Y-char.viewDist, 0)
		maxJ := min(char.Y+char.viewDist, len(maps[char.Mapname])-1)
		maxI := min(char.X+char.viewDist, len(maps[char.Mapname][0])-1)
		for ; j <= maxJ; j++ {
			i := max(char.X-char.viewDist, 0)
			for ; i <= maxI; i++ {
				if (j-char.Y)*(j-char.Y)+(i-char.X)*(i-char.X) > char.viewDist*char.viewDist+1 ||
					!visible(maps[char.Mapname], char.X, char.Y, i, j) {
					continue
				}
				if !hasMapPos(update.Maps[char.Mapname], i, j) {
					update.Maps[char.Mapname] = append(update.Maps[char.Mapname], mapPos{i, j, maps[char.Mapname][j][i]})
				}
			}
		}
	}
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func hasMapPos(pos []mapPos, x, y int) bool {
	for _, p := range pos {
		if p.X == x && p.Y == y {
			return true
		}
	}
	return false
}

func collectVisible(update *tickUpdate) {
	for mapname, pos := range update.Maps {
		for _, p := range pos {
			if char := findCharacter(mapname, p.X, p.Y); char != nil {
				update.Characters = append(update.Characters, char)
			}
			if prop, ok := findProp(mapname, p.X, p.Y); ok {
				update.Props = append(update.Props, prop)
			}
			if item, ok := findItem(mapname, p.X, p.Y); ok {
				update.Items = append(update.Items, item)
			}
		}
	}
}
func findCharacter(mapname string, x, y int) *Character {
	for _, char := range characters {
		if char.Mapname == mapname && char.X == x && char.Y == y {
			return char
		}
	}
	return nil
}
func findItem(mapname string, x, y int) (Item, bool) {
	for _, item := range items[mapname] {
		if item.X == x && item.Y == y {
			return item, true
		}
	}
	return Item{}, false
}
func findProp(mapname string, x, y int) (Prop, bool) {
	for _, prop := range props[mapname] {
		if prop.X == x && prop.Y == y {
			return prop, true
		}
	}
	return Prop{}, false
}
