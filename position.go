package main

type Position struct {
	Mapid string
	X     int
	Y     int
}

func (c *Character) getPosition() Position {
	return Position{c.Mapname, c.X, c.Y}
}

func (p *Position) getMapPos(m [][]bits) mapPos {
	return mapPos{
		Data: m[p.Y][p.X],
		X:    p.X,
		Y:    p.Y,
	}
}
