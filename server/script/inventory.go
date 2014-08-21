package script

import (
	"github.com/aarzilli/golua/lua"
	"github.com/elegios/topdown/server/types"
)

func (v *vm) give(L *lua.State) int {
	character := L.ToGoStruct(1).(*types.Character)
	item := L.CheckString(2)

	L.PushBoolean(character.AddItem(item))
	return 1
}

func (v *vm) take(L *lua.State) int {
	character := L.ToGoStruct(1).(*types.Character)
	item := L.CheckString(2)

	L.PushBoolean(character.RemoveItem(item))
	return 1
}
