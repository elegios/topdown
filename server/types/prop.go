package types

type Prop struct {
	X         int  `json:"x"`
	Y         int  `json:"y"`
	Variation int  `json:"variation"`
	Collide   bool `json:"collide"`
}

func (w *World) loadProps(path, mapname string) error {

	return nil
}
