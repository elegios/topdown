package types

type CharacterRunnable func(origin, target *Character)

const (
	EQUIP_NONE = iota
	EQUIP_WEAP
	EQUIP_ARMO
)

type Item struct {
	X  int    `json:"x"`
	Y  int    `json:"y"`
	Id string `json:"id"`
}

type ItemBlueprint struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Variation   int    `json:"variation"`
	Description string `json:"description"`

	Equippable int               `json:"-"`
	Effect     CharacterRunnable `json:"-"`
}

func (w *World) loadItems(path, mapname string) (err error) {
	var items []Item
	err = dec(path, &items)
	if err != nil {
		return
	}

	for _, item := range items {
		w.MapItems[Position{mapname, item.X, item.Y}] = item
	}

	return
}

func (w *World) saveItems(root string) (err error) {
	items := make(map[string][]Item)
	for pos, item := range w.MapItems {
		items[pos.Mapid] = append(items[pos.Mapid], item)
	}

	for name, itemlist := range items {
		err = enc(getPath(root, name, ITEM_EXT), itemlist)
		if err != nil {
			return
		}
	}

	return
}
