package script

import (
	"github.com/elegios/golua/lua"
)

type storyVm struct {
	vm
	autoend bool
	control chan struct{}
	module  string
	name    string
}

func (v *vm) RunStoryScript(module, name, path string, first bool) {
	s := storyVm{
		vm: vm{
			world: v.world,
			l:     lua.NewState(),
		},
		autoend: true,
		control: make(chan struct{}),
		module:  module,
		name:    name,
	}
	go s.run(path, first)
	<-s.control
}

func (s *storyVm) run(path string, first bool) {
	s.registerStory()

	d(s.trace(s.l.DoFile(path)))
	s.l.GetGlobal("main")
	s.l.PushBoolean(first)
	d(s.trace(s.l.Call(1, 0)))

	if s.autoend {
		s.remove()
	}
	s.control <- struct{}{}
}

func (s *storyVm) remove() {
	s.world.RemoveStory(s.module, s.name)
}
