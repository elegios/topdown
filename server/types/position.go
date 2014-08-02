package types

type Position struct {
	Mapid string `json:"mapname"`
	X     int    `json:"x"`
	Y     int    `json:"y"`
}

type MapPos struct {
	X    int  `json:"x"`
	Y    int  `json:"y"`
	Data Bits `json:"data"`
}

func (p *Position) GetMapPos(m [][]Bits) MapPos {
	return MapPos{
		Data: m[p.Y][p.X],
		X:    p.X,
		Y:    p.Y,
	}
}
