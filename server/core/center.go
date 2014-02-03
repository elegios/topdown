package core

import (
	"github.com/elegios/topdown/server/helpers"
	"github.com/elegios/topdown/server/types"
)

func center(ch <-chan request) {
	for req := range ch {
		c := req.message
		slog.Println("Got a command:", c)
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
		"name":        world.ItemBlueprints[id].Name,
		"type":        world.ItemBlueprints[id].Type,
		"variation":   world.ItemBlueprints[id].Variation,
		"description": world.ItemBlueprints[id].Description,
	}
}

func create(ch chan map[string]string, name string) {
	c := defaultCharacter
	c.Id = helpers.NewId(func(id string) bool {
		_, ok := world.Charids[id]
		return ok
	})
	c.Name = name
	world.Charids[c.Id] = &c
	world.MapCharacters[c.GetPosition()] = &c
	ch <- map[string]string{"id": c.Id}
}

func pickup(charId string) {
	c, ok := world.Charids[charId]
	if !ok || c.Actions < 1 {
		return
	}
	item, ok := world.MapItems[c.GetPosition()]
	if !ok {
		return
	}
	delete(world.MapItems, c.GetPosition())
	c.Inventory = append(c.Inventory, item.Id)
	c.Actions--
}

func useItem(charId, action, blueprintId string) {
	c, ok := world.Charids[charId]
	if !ok || c.Actions < 1 {
		return
	}
	index := -1
	for i, id := range c.Inventory {
		if id == blueprintId {
			index = i
			break
		}
	}
	if index == -1 {
		return
	}
	c.Actions--
	//TODO
	return
}

func move(charId, direction string) {
	c, ok := world.Charids[charId]
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
	if yMod < 0 || xMod < 0 || yMod >= len(world.Maps[c.Mapname]) || xMod >= len(world.Maps[c.Mapname][yMod]) {
		return
	}
	if world.Maps[c.Mapname][yMod][xMod].Collides() {
		return
	}
	pos := types.Position{c.Mapname, xMod, yMod}
	if _, ok := world.MapCharacters[pos]; ok {
		//TODO: attack
		return
	}
	if p, ok := world.MapProps[pos]; ok && p.Collide {
		//TODO: check if something special should happen
		return
	}
	c.Actions--
	delete(world.MapCharacters, c.GetPosition())
	c.X = xMod
	c.Y = yMod
	world.MapCharacters[c.GetPosition()] = c
}

func tick() {
	for _, c := range world.Charids {
		if c.Health <= 0 {
			delete(world.Charids, c.Id)
			delete(world.MapCharacters, c.GetPosition())
			continue
		}

		c.Actions = defaultActionCount
	}

	otm.In <- struct{}{}
	otm.Add(1)
	otm.Wait()
	slog.Print("Tick is done")
}
