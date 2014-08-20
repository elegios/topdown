package script

import (
	"github.com/elegios/golua/lua"
	"github.com/elegios/topdown/server/types"
)

func (v *vm) say(L *lua.State) int {
	if L.IsTable(2) {
		v.world.SpeakAt(checkPosition(L, 2), L.CheckString(1))

	} else {
		character := L.ToGoStruct(2).(*types.Character)
		v.world.Speak(character, L.CheckString(1))
	}

	return 0
}

func (v *vm) announce(L *lua.State) int {
	v.world.Announce(L.CheckString(1), L.CheckString(2))
	return 0
}
