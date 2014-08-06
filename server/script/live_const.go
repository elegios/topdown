package script

import (
	"github.com/aarzilli/golua/lua"
)

func (v *vm) registerConst(name string) {
	v.l.Register("weapon", v.weapon)
	v.l.Register("armor", v.armor)
	v.l.Register("item", v.itemConst)
	v.l.Register("gate", v.gate)
	v.l.PushString(name)
	v.l.SetField(lua.LUA_REGISTRYINDEX, "map")
}

func (v *vm) unregisterConst() {
	v.l.PushNil()
	v.l.SetGlobal("weapon")
	v.l.PushNil()
	v.l.SetGlobal("armor")
	v.l.PushNil()
	v.l.SetGlobal("item")
	v.l.PushNil()
	v.l.SetGlobal("gate")
}

func (v *vm) registerLive() {
	v.l.Register("item", v.itemLive)
	v.l.Register("nudge", v.nudge)
	v.l.Register("give", v.give)
	v.l.Register("take", v.take)
	v.l.Register("say", v.say)
}

func (v *vm) unregisterLive() {
	v.l.PushNil()
	v.l.SetGlobal("item")
	v.l.PushNil()
	v.l.SetGlobal("nudge")
	v.l.PushNil()
	v.l.SetGlobal("give")
	v.l.PushNil()
	v.l.SetGlobal("take")
	v.l.PushNil()
	v.l.SetGlobal("say")
}

func (v *vm) toggleToLive() {
	v.unregisterConst()
	v.registerLive()
}

func (v *vm) toggleToConst(name string) {
	v.unregisterLive()
	v.registerConst(name)
}
