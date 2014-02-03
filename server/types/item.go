package types

type Item struct {
	X  int    `json:"x"`
	Y  int    `json:"y"`
	Id string `json:"id"`
}

type ItemBlueprint struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Variation   string `json:"variation"`
	Description string `json:"description"`
}
