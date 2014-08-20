package script

import (
	"github.com/elegios/golua/lua"
	"github.com/elegios/topdown/server/types"
)

func (v *vm) gate(L *lua.State) int {
	L.GetField(lua.LUA_REGISTRYINDEX, "map")
	origin := types.Position{
		Mapid: L.ToString(-1),
		X:     L.CheckInteger(1),
		Y:     L.CheckInteger(2),
	}

	L.CheckType(3, lua.LUA_TTABLE)
	L.GetField(3, "map")
	L.GetField(3, "x")
	L.GetField(3, "y")

	v.world.MapTransitions[origin] = types.Position{
		Mapid: L.ToString(-3),
		X:     L.ToInteger(-2),
		Y:     L.ToInteger(-1),
	}

	return 0
}
