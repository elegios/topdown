package types

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

	Equippable int          `json:"-"`
	Effect     itemRunnable `json:"-"`
}
type itemRunnable func(origin, target *Character)

func (s *saved) loadItems(path, mapname string) (err error) {
	var items []Item
	err = dec(path, &items)
	if err != nil {
		return
	}

	for _, item := range items {
		s.MapItems[Position{mapname, item.X, item.Y}] = item
	}

	return
}

func (s *saved) saveItems(root string) (err error) {
	items := make(map[string][]Item)
	for pos, item := range s.MapItems {
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
