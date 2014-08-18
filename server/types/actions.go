package types

func (w *World) MoveCharacter(c *Character, direction string) bool {
	c.Actions--
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
		if pos, ok := w.MapTransitions[c.Pos]; ok {
			delete(w.MapCharacters, c.Pos)
			c.Pos = pos
			w.MapCharacters[c.Pos] = c
			return true
		}
		return false
	}
	if yMod < 0 || xMod < 0 || yMod >= len(w.Maps[c.Pos.Mapid]) || xMod >= len(w.Maps[c.Pos.Mapid][yMod]) {
		return false
	}
	if w.Maps[c.Pos.Mapid][yMod][xMod].Collides() {
		return false
	}
	pos := Position{c.Pos.Mapid, xMod, yMod}

	if other, ok := w.MapCharacters[pos]; ok {
		if c.Weapon != "" {
			w.ItemBlueprints[c.Weapon].Effect(c, other)
			return true
		}
		return false
	}

	p, ok := w.MapProps[pos]
	if ok && p.Collide {
		if p.Effect != nil {
			p.Effect(c)
			return true
		}
		return false
	}

	delete(w.MapCharacters, c.Pos)
	c.Pos.X = xMod
	c.Pos.Y = yMod
	w.MapCharacters[c.Pos] = c

	if ok && p.Effect != nil {
		p.Effect(c)
	}
	return true
}

func (w *World) UseItem(c *Character, blueprintId string) bool {
	index := -1
	for i, id := range c.Inventory {
		if id == blueprintId {
			index = i
			break
		}
	}
	if index == -1 {
		return false
	}
	blueprint, ok := w.ItemBlueprints[blueprintId]
	if !ok {
		return false
	}

	if blueprint.Equippable != EQUIP_NONE {
		c.RemoveItem(blueprintId)
		switch blueprint.Equippable {
		case EQUIP_WEAP:
			if c.Weapon != "" {
				c.AddItem(c.Weapon)
			}
			c.Weapon = blueprintId

		case EQUIP_ARMO:
			if c.Armor != "" {
				c.AddItem(c.Armor)
			}
			c.Armor = blueprintId
		}

	} else {
		blueprint.Effect(c, nil)
	}

	c.Actions--
	return true
}

func (w *World) DropItem(c *Character, blueprintId string) bool {
	index := -1
	for i, id := range c.Inventory {
		if id == blueprintId {
			index = i
			break
		}
	}
	if index == -1 {
		return false
	}
	pos := c.Pos
	if _, alreadyPresent := w.MapItems[pos]; alreadyPresent {
		return false
	}
	c.RemoveItem(blueprintId)
	w.MapItems[pos] = Item{c.Pos.X, c.Pos.Y, blueprintId}
	c.Actions--
	return true
}

func (w *World) PickupItem(c *Character) bool {
	item, ok := w.MapItems[c.Pos]
	if !ok {
		return false
	}
	delete(w.MapItems, c.Pos)
	c.Inventory = append(c.Inventory, item.Id)
	c.Actions--
	return true
}

func (w *World) Speak(c *Character, speech string) {
	p := c.Pos
	update := Update{
		Pos: &p,
		Content: SpeechCharUpdate{
			Speech:    speech,
			Character: c.Id,
		},
	}
	w.Updates = append(w.Updates, update)
}

func (w *World) SpeakAt(p Position, speech string) {
	update := Update{
		Pos: &p,
		Content: SpeechPointUpdate{
			Speech: speech,
			Pos:    p,
		},
	}
	w.Updates = append(w.Updates, update)
}
