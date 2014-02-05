package types

type Character struct {
	Mapname   string   `json:"mapname"`
	X         int      `json:"x"`
	Y         int      `json:"y"`
	Id        string   `json:"id"`
	Variation int      `json:"variation"`
	Name      string   `json:"name"`
	Actions   int      `json:"actions"`
	Weapon    string   `json:"weapon"`
	Armor     string   `json:"armor"`
	Health    int      `json:"health"`
	MaxHealth int      `json:"maxhealth"`
	Inventory []string `json:"inventory"`

	ViewDist int `json:"-"`
}

func (c *Character) GetPosition() Position {
	return Position{c.Mapname, c.X, c.Y}
}

func (w *World) loadCharacters(path string) (err error) {
	err = dec(path, &w.Charids)
	if err != nil {
		return
	}

	for _, c := range w.Charids {
		w.MapCharacters[c.GetPosition()] = c
	}

	return
}

func (w *World) saveCharacters(path string) error {
	return enc(path, w.Charids)
}
