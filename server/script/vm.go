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

func (v *VM) RunConstantScript(path string, world *types.World) (err error) {
	vm := (*gelo.VM)(v)
	vm.RegisterBundle(map[string]interface{}{
		"itemb": (*vmworld)(world).ItemBlueprint,
	})
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
