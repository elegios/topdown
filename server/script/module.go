package script

import (
	"github.com/aarzilli/golua/lua"
)

func (v *vm) applyModule(L *lua.State) int {
	v.world.ApplyModule(L.CheckString(1))

	return 0
}
