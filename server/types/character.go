package types

import (
	"github.com/elegios/topdown/server/helpers"
)

type Character struct {
	Pos           Position `json:"position"`
	Id            string   `json:"id"`
	Variation     int      `json:"variation"`
	Name          string   `json:"name"`
	Actions       float32  `json:"actions"`
	RecoverySpeed float32  `json:"-"`
	RecoveryMax   int      `json:"-"`
	Weapon        string   `json:"weapon"`
	Armor         string   `json:"armor"`
	Health        int      `json:"health"`
	MaxHealth     int      `json:"maxhealth"`
	Inventory     []string `json:"inventory"`

	ViewDist int `json:"-"`
}

func (c *Character) AddItem(bid string) bool {
	c.Inventory = append(c.Inventory, bid)
	return true //This will eventually be some check concerning inventory size
}

func (c *Character) RemoveItem(bid string) bool {
	for i, item := range c.Inventory {
		if item == bid {
			c.Inventory[i] = c.Inventory[len(c.Inventory)-1]
			c.Inventory = c.Inventory[:len(c.Inventory)-1]
			return true
		}
	}

	return false
}

func (w *World) loadCharacters(path string) (err error) {
	err = dec(path, &w.Charids)
	if err != nil {
		return
	}

	for _, c := range w.Charids {
		w.MapCharacters[c.Pos] = c
	}

	return
}

func (w *World) saveCharacters(path string) error {
	return enc(path, w.Charids)
}

func (c *Character) UpdateActions() {
	c.Actions = helpers.Maxf(
		c.Actions,
		helpers.Minf(
			float32(c.RecoveryMax),
			c.Actions+c.RecoverySpeed))
}
