package script

import (
	"github.com/aarzilli/golua/lua"
	"github.com/elegios/topdown/server/types"
)

func (v *vm) say(L *lua.State) int {
	if L.IsTable(2) {
		L.GetField(2, "map")
		L.GetField(2, "x")
		L.GetField(2, "y")
		pos := types.Position{
			Mapid: L.ToString(-3),
			X:     L.ToInteger(-2),
			Y:     L.ToInteger(-1),
		}
		update := types.Update{
			Pos: pos,
			Content: types.SpeechPointUpdate{
				Pos:    pos,
				Speech: L.CheckString(1),
			},
		}
		v.world.Updates = append(v.world.Updates, update)

	} else {
		character := L.ToGoStruct(2).(*types.Character)
		update := types.Update{
			Pos: character.Pos,
			Content: types.SpeechCharUpdate{
				Character: character.Id,
				Speech:    L.CheckString(1),
			},
		}
		v.world.Updates = append(v.world.Updates, update)
	}

	return 0
}
