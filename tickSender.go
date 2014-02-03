package main

import (
	"code.google.com/p/go.net/websocket"
	"log"
)

type tickUpdate struct {
	Maps         map[string][]mapPos `json:"maps"`
	Controllable []string            `json:"controllable"`
	Characters   []*Character        `json:"characters"`
	Props        []Prop              `json:"props"`
	Items        []Item              `json:"items"`
	positions    map[Position]bool
}

type mapPos struct {
	X    int  `json:"x"`
	Y    int  `json:"y"`
	Data bits `json:"data"`
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
			positions:    make(map[Position]bool),
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
		_, ok := world.charids[charId]
		if !ok {
			delete(data.charIds, charId)
			continue
		}
		update.Controllable = append(update.Controllable, charId)
	}
}

func collectMapPos(update *tickUpdate) {
	for _, cid := range update.Controllable {
		char := world.charids[cid]
		j := max(char.Y-char.viewDist, 0)
		maxJ := min(char.Y+char.viewDist, len(world.maps[char.Mapname])-1)
		maxI := min(char.X+char.viewDist, len(world.maps[char.Mapname][0])-1)
		pos := Position{
			mapid: char.Mapname,
		}
		for ; j <= maxJ; j++ {
			pos.y = j
			i := max(char.X-char.viewDist, 0)
			for ; i <= maxI; i++ {
				pos.x = i
				if update.positions[pos] ||
					(j-char.Y)*(j-char.Y)+(i-char.X)*(i-char.X) > char.viewDist*char.viewDist+1 ||
					!visible(world.maps[char.Mapname], char.X, char.Y, i, j) {
					continue
				}
				update.positions[Position{char.Mapname, i, j}] = true
			}
		}
	}
	for pos := range update.positions {
		update.Maps[pos.mapid] = append(update.Maps[pos.mapid], mapPos{pos.x, pos.y, world.maps[pos.mapid][pos.y][pos.x]})
	}
}

func collectVisible(update *tickUpdate) {
	for pos := range update.positions {
		if c, ok := world.mapCharacters[pos]; ok {
			update.Characters = append(update.Characters, c)
		}
		if p, ok := world.mapProps[pos]; ok {
			update.Props = append(update.Props, p)
		}
		if i, ok := world.mapItems[pos]; ok {
			update.Items = append(update.Items, i)
		}
	}
}
