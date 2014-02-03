package main

func (c *Character) getPosition() Position {
	return Position{c.Mapname, c.X, c.Y}
}
