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

func blueprintRequest(ch chan map[string]interface{}, id string) {
	ch <- map[string]interface{}{
		"id":          id,
		"name":        world.ItemBlueprints[id].Name,
		"type":        world.ItemBlueprints[id].Type,
		"variation":   world.ItemBlueprints[id].Variation,
		"description": world.ItemBlueprints[id].Description,
	}
}

func create(ch chan map[string]interface{}, name string) {
	c := defaultCharacter
	c.Id = helpers.NewId(func(id string) bool {
		_, ok := world.Charids[id]
		return ok
	})
	c.Name = name
	c.Inventory = []string{"hpotmed", "hpotmin", "cream", "stabbysword"}
	world.Charids[c.Id] = &c
	world.MapCharacters[c.Pos] = &c
	ch <- map[string]interface{}{"id": c.Id}
}

func pickup(charId string) {
	c, ok := world.Charids[charId]
	if !ok || c.Actions < 1 {
		return
	}
	item, ok := world.MapItems[c.Pos]
	if !ok {
		return
	}
	delete(world.MapItems, c.Pos)
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
	switch action {
	case "use":
		blueprint, ok := world.ItemBlueprints[blueprintId]
		if !ok {
			return
		}

		if blueprint.Equippable != types.EQUIP_NONE {
			c.RemoveItem(blueprintId)
			switch blueprint.Equippable {
			case types.EQUIP_WEAP:
				if c.Weapon != "" {
					c.AddItem(c.Weapon)
				}
				c.Weapon = blueprintId

			case types.EQUIP_ARMO:
				if c.Armor != "" {
					c.AddItem(c.Armor)
				}
				c.Armor = blueprintId
			}

		} else {
			world.ItemBlueprints[blueprintId].Effect.RunOn(c, nil)
		}

	case "drop":
		pos := c.Pos
		if _, alreadyPresent := world.MapItems[pos]; alreadyPresent {
			return
		}
		c.RemoveItem(blueprintId)
		world.MapItems[pos] = types.Item{c.Pos.X, c.Pos.Y, blueprintId}
	}
	return
}

func move(charId, direction string) {
	c, ok := world.Charids[charId]
	if !ok || c.Actions < 1 {
		return
	}
	xMod := c.Pos.X
	yMod := c.Pos.Y
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
		if pos, ok := world.MapTransitions[c.Pos]; ok {
			delete(world.MapCharacters, c.Pos)
			c.Pos = pos
			world.MapCharacters[c.Pos] = c
		}
		wlog.Print(c.Pos)
		return
	}
	if yMod < 0 || xMod < 0 || yMod >= len(world.Maps[c.Pos.Mapid]) || xMod >= len(world.Maps[c.Pos.Mapid][yMod]) {
		return
	}
	if world.Maps[c.Pos.Mapid][yMod][xMod].Collides() {
		return
	}
	pos := types.Position{c.Pos.Mapid, xMod, yMod}
	if _, ok := world.MapCharacters[pos]; ok {
		//TODO: attack
		return
	}
	if p, ok := world.MapProps[pos]; ok && p.Collide {
		//TODO: check if something special should happen
		return
	}
	c.Actions--
	delete(world.MapCharacters, c.Pos)
	c.Pos.X = xMod
	c.Pos.Y = yMod
	world.MapCharacters[c.Pos] = c
}

func tick() {
	for _, c := range world.Charids {
		if c.Health <= 0 {
			delete(world.Charids, c.Id)
			delete(world.MapCharacters, c.Pos)
			continue
		}

		c.Actions = defaultActionCount
	}

	otm.In <- struct{}{}
	otm.Add(1)
	otm.Wait()
	world.Updates = nil
	slog.Print("Tick is done")
}
