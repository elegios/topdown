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
		case "blueprintrequest":
			blueprintRequest(req.ch, c["id"])

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

func blueprintRequest(ch chan map[string]string, id string) {
	ch <- map[string]string{
		"id":          id,
		"name":        world.itemBlueprints[id].Name,
		"type":        world.itemBlueprints[id].Type,
		"variation":   world.itemBlueprints[id].Variation,
		"description": world.itemBlueprints[id].Description,
	}
}

func create(ch chan map[string]string, name string) {
	b := make([]byte, idlength)
	var id string
	for inUse := true; inUse; _, inUse = world.charids[id] {
		rand.Read(b)
		id = base64.StdEncoding.EncodeToString(b)
	}
	c := defaultCharacter
	c.Id = id
	c.Name = name
	world.charids[id] = &c
	world.mapCharacters[c.getPosition()] = &c
	ch <- map[string]string{"id": id}
}

func pickup(charId string) {
	c, ok := world.charids[charId]
	if !ok || c.Actions < 1 {
		return
	}
	item, ok := world.mapItems[c.getPosition()]
	if !ok {
		return
	}
	delete(world.mapItems, c.getPosition())
	c.Inventory = append(c.Inventory, item.Id)
	c.Actions--
}

func useItem(charId, action, itemId string) {
	//TODO
}

func move(charId, direction string) {
	c, ok := world.charids[charId]
	if !ok || c.Actions < 1 {
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
	if yMod < 0 || xMod < 0 || yMod >= len(world.maps[c.Mapname]) || xMod >= len(world.maps[c.Mapname][yMod]) {
		return
	}
	if world.maps[c.Mapname][yMod][xMod].collides() {
		return
	}
	pos := Position{c.Mapname, xMod, yMod}
	if _, ok := world.mapCharacters[pos]; ok {
		//TODO: attack
		return
	}
	if p, ok := world.mapProps[pos]; ok && p.Collide {
		//TODO: check if something special should happen
		return
	}
	c.Actions--
	delete(world.mapCharacters, c.getPosition())
	c.X = xMod
	c.Y = yMod
	world.mapCharacters[c.getPosition()] = c
}

func tick() {
	for _, c := range world.charids {
		if c.Health <= 0 {
			delete(world.charids, c.Id)
			delete(world.mapCharacters, c.getPosition())
			continue
		}

		c.Actions = defaultActionCount
	}

	otm.In <- struct{}{}
	otm.Add(1)
	otm.Wait()
	log.Print("Tick is done")
}
