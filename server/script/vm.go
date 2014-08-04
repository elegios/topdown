package script

import (
	"code.google.com/p/gelo"
	"github.com/elegios/topdown/server/types"
	"os"
)

type VM gelo.VM
type constArgs struct {
	world *types.World
	name  string
}
type vmworld types.World

func CreateVM() *VM {
	v := gelo.NewVM(nil)
	v.RegisterBundles(langCommands)
	v.Ns.Fork(nil)
	return (*VM)(v)
}

func (v *VM) RunConstantScript(path, name string, world *types.World) error {
	w := &constArgs{
		world: world,
		name:  name,
	}
	return v.runScript(path, constBundle(w))
}
func constBundle(world *constArgs) map[string]interface{} {
	return map[string]interface{}{
		"itemb": world.ItemBlueprint,
		"gate":  world.Gate,
	}
}

func (v *VM) RunLiveScript(path string, world *types.World) error {
	return v.runScript(path, liveBundle(world))
}
func liveBundle(world *types.World) map[string]interface{} {
	return map[string]interface{}{
		"item":  (*vmworld)(world).Item,
		"speak": (*vmworld)(world).SpeakAt,
	}
}

func (v *VM) runScript(path string, bundle map[string]interface{}) (err error) {
	vm := (*gelo.VM)(v)
	vm.Ns.Fork(nil)
	defer vm.Ns.Unfork()
	vm.RegisterBundle(bundle)
	vm.Ns.Fork(nil)
	defer vm.Ns.Unfork()
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	_, err = vm.Run(f, nil)
	return
}
