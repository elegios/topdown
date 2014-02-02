package main

import (
	"code.google.com/p/go.net/websocket"
	"log"
)

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
		_, ok := world.characters[charId]
		if !ok {
			delete(data.charIds, charId)
			continue
		}
		update.Controllable = append(update.Controllable, charId)
	}
}

func collectMapPos(update *tickUpdate) {
	for _, cid := range update.Controllable {
		char := world.characters[cid]
		j := max(char.Y-char.viewDist, 0)
		maxJ := min(char.Y+char.viewDist, len(world.maps[char.Mapname])-1)
		maxI := min(char.X+char.viewDist, len(world.maps[char.Mapname][0])-1)
		for ; j <= maxJ; j++ {
			i := max(char.X-char.viewDist, 0)
			for ; i <= maxI; i++ {
				if (j-char.Y)*(j-char.Y)+(i-char.X)*(i-char.X) > char.viewDist*char.viewDist+1 ||
					!visible(world.maps[char.Mapname], char.X, char.Y, i, j) {
					continue
				}
				if !hasMapPos(update.Maps[char.Mapname], i, j) {
					update.Maps[char.Mapname] = append(update.Maps[char.Mapname], mapPos{i, j, world.maps[char.Mapname][j][i]})
				}
			}
		}
	}
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
	for _, char := range world.characters {
		if char.Mapname == mapname && char.X == x && char.Y == y {
			return char
		}
	}
	return nil
}
func findItem(mapname string, x, y int) (Item, bool) {
	for _, item := range world.items[mapname] {
		if item.X == x && item.Y == y {
			return item, true
		}
	}
	return Item{}, false
}
func findProp(mapname string, x, y int) (Prop, bool) {
	for _, prop := range world.props[mapname] {
		if prop.X == x && prop.Y == y {
			return prop, true
		}
	}
	return Prop{}, false
}
