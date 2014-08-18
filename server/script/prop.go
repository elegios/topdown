package script

import (
	"github.com/aarzilli/golua/lua"
	"github.com/elegios/topdown/server/types"
)

func (v *vm) prop(L *lua.State) int {
	p := checkPosition(L, 2)

	if L.IsNil(1) {
		delete(v.world.MapProps, p)
		return 0
	}

	L.CheckType(1, lua.LUA_TTABLE)
	L.GetField(1, "variation")
	L.GetField(1, "collide")

	v.world.MapProps[p] = types.Prop{
		X:         p.X,
		Y:         p.Y,
		Variation: L.ToInteger(-2),
		Collide:   L.ToBoolean(-1),
	}
	return 0
}

func (s *storyVm) prop(L *lua.State) int {
	p := checkPosition(L, 2)

	if L.IsNil(1) {
		delete(s.world.MapProps, p)
		return 0
	}

	L.CheckType(1, lua.LUA_TTABLE)
	L.GetField(1, "variation")
	L.GetField(1, "collide")
	var thing interface{} = struct{}{}
	L.PushLightUserdata(&thing)
	L.GetField(1, "effect")

	if L.IsNil(-1) {
		s.world.MapProps[p] = types.Prop{
			X:         p.X,
			Y:         p.Y,
			Variation: L.ToInteger(3),
			Collide:   L.ToBoolean(4),
		}

	} else {
		L.SetTable(lua.LUA_REGISTRYINDEX)
		s.world.MapProps[p] = types.Prop{
			X:         p.X,
			Y:         p.Y,
			Variation: L.ToInteger(3),
			Collide:   L.ToBoolean(4),
			Effect: func(c *types.Character) {
				go func() {
					L.PushLightUserdata(&thing)
					L.GetTable(lua.LUA_REGISTRYINDEX)
					L.PushGoStruct(c)
					d(s.trace(L.Call(1, 0)))
					s.control <- struct{}{}
				}()
				<-s.control
			},
		}
	}

	return 0
}
