package script

import (
	"code.google.com/p/gelo"
	"github.com/elegios/topdown/server/types"
	"os"
)

type VM gelo.VM
type vmworld types.World

func CreateVM() *VM {
	v := gelo.NewVM(nil)
	v.RegisterBundles(langCommands)
	v.Ns.Fork(nil)
	return (*VM)(v)
}

func (v *VM) RunConstantScript(path string, world *types.World) error {
	return v.runScript(path, world, map[string]interface{}{
		"itemb": (*vmworld)(world).ItemBlueprint,
	})
}

func (v *VM) RunLiveScript(path string, world *types.World) error {
	return v.runScript(path, world, map[string]interface{}{
		"item": (*vmworld)(world).Item,
	})
}

func (v *VM) runScript(path string, world *types.World, bundle map[string]interface{}) (err error) {
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
