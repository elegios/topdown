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

		case "speak":
			speak(c["character"], c["speech"])
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
	c.Inventory = []string{"hpotmed", "hpotmin", "cream", "stabbysword", "speaker"}
	world.Charids[c.Id] = &c
	world.MapCharacters[c.Pos] = &c
	ch <- map[string]interface{}{"id": c.Id}
}

func pickup(charId string) {
	c, ok := world.Charids[charId]
	if !ok || c.Actions < 1 {
		return
	}
	world.PickupItem(c)
}

func useItem(charId, action, blueprintId string) {
	c, ok := world.Charids[charId]
	if !ok || c.Actions < 1 {
		return
	}
	switch action {
	case "use":
		world.UseItem(c, blueprintId)
	case "drop":
		world.DropItem(c, blueprintId)
	}
}

func move(charId, direction string) {
	c, ok := world.Charids[charId]
	if !ok || c.Actions < 1 {
		return
	}
	world.MoveCharacter(c, direction)
}

func tick() {
	for _, c := range world.Charids {
		if c.Health <= 0 {
			delete(world.Charids, c.Id)
			delete(world.MapCharacters, c.Pos)
			continue
		}

		c.UpdateActions()
	}

	otm.In <- struct{}{}
	otm.Add(1)
	otm.Wait()
	world.Updates = nil
	slog.Print("Tick is done")
}

func speak(cid, speech string) {
	c, ok := world.Charids[cid]
	if !ok {
		return
	}

	world.Speak(c, speech)
}
