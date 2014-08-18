package script

import (
	"github.com/aarzilli/golua/lua"
	"github.com/elegios/topdown/server/types"
	"runtime"
)

func (s *storyVm) autoEnd(L *lua.State) int {
	s.autoend = L.ToBoolean(1)
	return 0
}

func (s *storyVm) endStory(L *lua.State) int {
	s.remove()
	s.control <- struct{}{}
	L.Close()
	runtime.Goexit()
	return 0 // Compiler pleasing, Goexit will stop execution
}

func (s *storyVm) retrieveValue(L *lua.State) int {
	v := s.world.ReadStoryValue(s.module, s.name, L.CheckString(1))

	switch v := v.(type) {
	case bool:
		L.PushBoolean(v)

	case float32:
		L.PushNumber(float64(v))

	case string:
		L.PushString(v)

	case types.Position:
		L.CreateTable(0, 3)
		L.PushString(v.Mapid)
		L.SetField(2, "map")
		L.PushInteger(int64(v.X))
		L.SetField(2, "x")
		L.PushInteger(int64(v.Y))
		L.SetField(2, "y")

	default:
		panic("Unsupported type retrieved")
	}

	return 1
}

func (s *storyVm) storeValue(L *lua.State) int {
	k := L.CheckString(1)
	var v interface{}

	switch {
	case L.IsBoolean(2):
		v = L.ToBoolean(2)

	case L.IsNumber(2):
		v = float32(L.ToNumber(2))

	case L.IsString(2):
		v = L.ToString(2)

	case L.IsTable(2):
		v = checkPosition(L, 2)

	default:
		panic("Unsupported type stored")
	}

	s.world.SetStoryValue(s.module, s.name, k, v)
	return 0
}

func (s *storyVm) applyPartial(L *lua.State) int {
	d(s.world.ApplyPartial(checkPosition(L, 2), s.module, L.CheckString(1)))

	return 0
}

func (s *storyVm) subStory(L *lua.State) int {
	s.world.RunStory(s.module, L.CheckString(1))
	return 0
}
