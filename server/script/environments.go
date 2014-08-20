package script

import (
	"github.com/elegios/golua/lua"
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
	v.l.Register("prop", v.prop)
	v.l.Register("announce", v.announce)
	v.l.Register("apply_module", v.applyModule)
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
	v.l.PushNil()
	v.l.SetGlobal("prop")
	v.l.PushNil()
	v.l.SetGlobal("announce")
	v.l.PushNil()
	v.l.SetGlobal("apply_module")
}

func (s *storyVm) registerStory() {
	s.registerLive()
	s.l.Register("prop", s.prop)
	s.l.Register("retrieve_value", s.retrieveValue)
	s.l.Register("store_value", s.storeValue)
	s.l.Register("auto_end", s.autoEnd)
	s.l.Register("apply_partial", s.applyPartial)
	s.l.Register("substory", s.subStory)
	s.l.Register("end_story", s.endStory)
}

func (v *vm) toggleToLive() {
	v.unregisterConst()
	v.registerLive()
}

func (v *vm) toggleToConst(name string) {
	v.unregisterLive()
	v.registerConst(name)
}
