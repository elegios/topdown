package main

import (
	"crypto/rand"
	"encoding/base64"
	"log"
)

func Center(ch <-chan request) {
	for req := range ch {
		c := req.message
		log.Println("Got a command:", c)
		switch c["command"] {
		case "itemrequest":
			itemRequest(req.ch, c["id"])

		case "create":
			create(req.ch, c["name"])

		case "pickup":
			pickup(c["character"])

		case "useitem":
			useItem(c["character"], c["action"], c["item"])

		case "move":
			move(c["character"], c["direction"])

		case "tick":
			tick()
		}
	}
}

func itemRequest(ch chan map[string]string, id string) {
	ch <- map[string]string{
		"id":          itemData[id].Id,
		"name":        itemData[id].Name,
		"type":        itemData[id].Type,
		"variation":   itemData[id].Variation,
		"description": itemData[id].Description,
	}
}

func create(ch chan map[string]string, name string) {
	b := make([]byte, idlength)
	var id string
	for inUse := true; inUse; _, inUse = characters[id] {
		rand.Read(b)
		id = base64.StdEncoding.EncodeToString(b)
	}
	c := defaultCharacter
	c.Id = id
	c.Name = name
	characters[id] = &c
	ch <- map[string]string{"id": id}
}

func pickup(charId string) {
	c := characters[charId]
	if c == nil || c.Actions < 1 {
		return
	}
	for i, item := range items[c.Mapname] {
		if item.X == c.X && item.Y == c.Y {
			items[c.Mapname][i] = items[c.Mapname][len(items[c.Mapname])-1]
			items[c.Mapname] = items[c.Mapname][:len(items[c.Mapname])-1]
			c.Inventory = append(c.Inventory, item.Id)
			c.Actions--
			return
		}
	}
}

func useItem(charId, action, itemId string) {
	//TODO
}

func move(charId, direction string) {
	c := characters[charId]
	if c == nil || c.Actions < 1 {
		return
	}
	xMod := c.X
	yMod := c.Y
	switch direction {
	case "left":
		xMod += -1
	case "right":
		xMod += 1
	case "up":
		yMod += -1
	case "down":
		yMod += 1
	case "through":
		//TODO: changing maps
		return
	}
	if yMod < 0 || xMod < 0 || yMod >= len(maps[c.Mapname]) || xMod >= len(maps[c.Mapname][yMod]) {
		return
	}
	if maps[c.Mapname][yMod][xMod].collides() {
		return
	}
	for _, c2 := range characters {
		if c.Mapname == c2.Mapname && c2.X == xMod && c2.Y == yMod {
			//TODO: attack
			return
		}
	}
	for _, p := range props[c.Mapname] {
		if p.Collide && p.X == xMod && p.Y == yMod {
			//TODO: check if something special should happen
			return
		}
	}
	c.Actions--
	c.X = xMod
	c.Y = yMod
}

func tick() {
	for _, c := range characters {
		if c.Health <= 0 {
			delete(characters, c.Id)
			continue
		}

		c.Actions = defaultActionCount
	}

	otm.In <- struct{}{}
	otm.Add(1)
	otm.Wait()
	log.Print("Tick is done")
}
