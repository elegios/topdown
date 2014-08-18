package script

import (
	"github.com/aarzilli/golua/lua"
	"github.com/elegios/topdown/server/types"
)

func checkPosition(L *lua.State, arg int) types.Position {
	L.CheckType(arg, lua.LUA_TTABLE)
	L.GetField(arg, "map")
	L.GetField(arg, "x")
	L.GetField(arg, "y")
	defer L.Pop(3)

	return types.Position{
		Mapid: L.ToString(-3),
		X:     L.ToInteger(-2),
		Y:     L.ToInteger(-1),
	}
}
