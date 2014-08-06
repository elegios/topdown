package script

import (
	"github.com/aarzilli/golua/lua"
	"github.com/elegios/topdown/server/types"
)

func (v *vm) weapon(L *lua.State) int {
	return v.itemb(L, types.EQUIP_WEAP)
}
func (v *vm) armor(L *lua.State) int {
	return v.itemb(L, types.EQUIP_ARMO)
}
func (v *vm) itemConst(L *lua.State) int {
	return v.itemb(L, types.EQUIP_NONE)
}

func (v *vm) itemb(L *lua.State, equippable int) int {
	b := types.ItemBlueprint{
		Equippable: equippable,
	}

	L.CheckType(1, lua.LUA_TTABLE)

	L.GetField(1, "id")
	id := L.ToString(-1)
	runner := "ib." + id

	L.GetField(1, "name")
	b.Name = L.ToString(-1)

	L.GetField(1, "type")
	b.Type = L.ToString(-1)

	L.GetField(1, "variation")
	b.Variation = L.ToInteger(-1)

	L.GetField(1, "description")
	b.Description = L.ToString(-1)

	L.GetField(1, "effect")
	L.SetField(lua.LUA_REGISTRYINDEX, runner)
	b.Effect = func(origin, target *types.Character) {
		v.l.GetField(lua.LUA_REGISTRYINDEX, runner)
		v.l.PushGoStruct(origin)
		v.l.PushGoStruct(target)
		v.trace(v.l.Call(2, 0))
	}

	v.world.ItemBlueprints[id] = b

	return 0
}

func (v *vm) itemLive(L *lua.State) int {
	L.CheckType(2, lua.LUA_TTABLE)
	L.GetField(2, "map")
	L.GetField(2, "x")
	L.GetField(2, "y")
	pos := types.Position{
		Mapid: L.ToString(-3),
		X:     L.ToInteger(-2),
		Y:     L.ToInteger(-1),
	}
	v.world.MapItems[pos] = types.Item{ //TODO: fail when already exists?
		X:  pos.X,
		Y:  pos.Y,
		Id: L.CheckString(1),
	}
	return 0
}
