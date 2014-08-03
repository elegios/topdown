package script

import (
	"code.google.com/p/gelo"
	"code.google.com/p/gelo/extensions"
	"github.com/elegios/topdown/server/types"
)

var (
	gateParser = extensions.MakeOrElseArgParser("'at ox oy 'to name 'at tx ty")
)

func (c *constArgs) Gate(vm *gelo.VM, args *gelo.List, argc uint) gelo.Word {
	dict := gateParser(vm, args)
	ox := ensureInt(vm, dict["ox"])
	oy := ensureInt(vm, dict["oy"])
	name := vm.API.SymbolOrElse(dict["name"]).String()
	tx := ensureInt(vm, dict["tx"])
	ty := ensureInt(vm, dict["ty"])

	origin := types.Position{
		Mapid: c.name,
		X:     ox,
		Y:     oy,
	}
	target := types.Position{
		Mapid: name,
		X:     tx,
		Y:     ty,
	}
	c.world.MapTransitions[origin] = target

	return gelo.Null
}
