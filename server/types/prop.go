package types

type Prop struct {
	X         int  `json:"x"`
	Y         int  `json:"y"`
	Variation int  `json:"variation"`
	Collide   bool `json:"collide"`
}

func (w *World) loadProps(path, mapname string) (err error) {
	var props []Prop
	err = dec(path, &props)
	if err != nil {
		return
	}

	for _, prop := range props {
		w.MapProps[Position{mapname, prop.X, prop.Y}] = prop
	}

	return
}

func (w *World) saveProps(root string) (err error) {
	props := make(map[string][]Prop)
	for pos, prop := range w.MapProps {
		props[pos.Mapid] = append(props[pos.Mapid], prop)
	}

	for name, proplist := range props {
		err = enc(getPath(root, name, PROP_EXT), proplist)
		if err != nil {
			return
		}
	}

	return
}
