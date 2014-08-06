package script

import (
	"github.com/aarzilli/golua/lua"
	"github.com/elegios/topdown/server/types"
	"log"
)

type vm struct {
	world *types.World
	l     *lua.State
}

func (v *vm) initVM(L *lua.State) int {
	v.registerLive()
	return 0
}

func LoadWorld(root string) *types.World {
	v := &vm{
		world: new(types.World),
		l:     lua.NewState(),
	}
	v.l.PushGoFunction(v.initVM)
	v.l.Call(0, 0)

	types.LoadWorld(v.world, v, root)

	return v.world
}

func (v *vm) RunConstantScript(path, name string) error {
	v.toggleToConst(name)
	defer v.toggleToLive()
	return v.trace(v.l.DoFile(path))
}

func (v *vm) RunLiveScript(path string) error {
	return v.trace(v.l.DoFile(path))
}

func d(err error) {
	if err != nil {
		panic(err)
	}
}

func (v *vm) trace(err error) error {
	if err == nil {
		return nil
	}
	log.Println(err)
	log.Println(v.l.StackTrace())
	return err
}
