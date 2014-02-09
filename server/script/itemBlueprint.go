package script

import (
	"code.google.com/p/gelo"
	"github.com/elegios/topdown/server/types"
)

type itemRunner struct {
	world *types.World
	vm    *gelo.VM
	code  gelo.Quote
	ns    *gelo.Dict
}

func (w *vmworld) ItemBlueprint(vm *gelo.VM, args *gelo.List, argc uint) gelo.Word {
	if argc != 6 {
		gelo.ArgumentError(vm, "itemb", "string string string int string code", args)
	}

	iid := vm.API.SymbolOrElse(args.Value).String()
	args = args.Next
	iname := vm.API.SymbolOrElse(args.Value).String()
	args = args.Next
	itype := vm.API.SymbolOrElse(args.Value).String()
	args = args.Next
	ivariation, _ := vm.API.NumberOrElse(args.Value).Int()
	args = args.Next
	idescription := vm.API.SymbolOrElse(args.Value).String()
	args = args.Next
	icode := vm.API.QuoteOrElse(args.Value)

	w.ItemBlueprints[iid] = types.ItemBlueprint{
		Name:        iname,
		Type:        itype,
		Variation:   int(ivariation),
		Description: idescription,
		Effect: &itemRunner{
			world: (*types.World)(w),
			vm:    vm,
			code:  icode,
			ns:    vm.Ns.Locals(vm.Ns.Depth() - 3), // ignore the namespaces with the language and the "extraapi"
		},
	}

	return gelo.Null
}

func Nudge(char *types.Character, val string, amount int) {
	switch val {
	case "health":
		char.Health += amount

	case "maxhealth":
		char.MaxHealth += amount

	case "viewdist":
		char.ViewDist += amount
	}
}

// target may be nil, origin shouldn't be
func (i *itemRunner) RunOn(origin, target *types.Character) {
	i.vm.Ns.Fork(nil)
	defer i.vm.Ns.Unfork()
	i.vm.Ns.Inject(0, i.ns)

	vals := &onAttackVals{
		origin: origin,
		target: target,
		world:  i.world,
	}

	i.vm.RegisterBundle(map[string]interface{}{
		"nudge":  vals.Nudge,
		"remove": vals.RemoveItem,
	})

	i.vm.SetProgram(i.code)
	i.vm.Exec(nil)
}

type onAttackVals struct {
	origin, target *types.Character
	world          *types.World
}

func (o *onAttackVals) Nudge(vm *gelo.VM, args *gelo.List, argc uint) gelo.Word {
	if argc != 2 && argc != 3 {
		gelo.ArgumentError(vm, "nudge", "[me|target] string integer", args)
	}
	c := o.origin
	if argc == 3 {
		switch vm.API.SymbolOrElse(args.Value).String() {
		case "me":
			c = o.origin

		case "target":
			c = o.target
			if c == nil {
				gelo.RuntimeError(vm, "There is no target for this action.")
			}

		default:
			gelo.ArgumentError(vm, "nudge", "[me|target] string integer", args)
		}
		args = args.Next
	}
	val := vm.API.SymbolOrElse(args.Value).String()
	amount, _ := vm.API.NumberOrElse(args.Next.Value).Int()
	Nudge(c, val, int(amount))

	return gelo.Null
}

func (o *onAttackVals) RemoveItem(vm *gelo.VM, args *gelo.List, argc uint) gelo.Word {
	if argc != 1 && argc != 2 {
		gelo.ArgumentError(vm, "remove", "[me|target] string", args)
	}
	c := o.origin
	if argc == 2 {
		switch vm.API.SymbolOrElse(args.Value).String() {
		case "me":
			c = o.origin

		case "target":
			c = o.target
			if c == nil {
				gelo.RuntimeError(vm, "There is no target for this action.")
			}

		default:
			gelo.ArgumentError(vm, "remove", "[me|target] string", args)
		}
		args = args.Next
	}

	bid := vm.API.SymbolOrElse(args.Value).String()
	for i, item := range c.Inventory {
		if item == bid {
			c.Inventory[i] = c.Inventory[len(c.Inventory)-1]
			c.Inventory = c.Inventory[:len(c.Inventory)-1]
			return gelo.True
		}
	}
	return gelo.False
}
