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

func (w *vmworld) Item(vm *gelo.VM, args *gelo.List, argc uint) gelo.Word {
	if argc > 1 {
		gelo.ArgumentError(vm, "item", "[properties]", args)
	}

	var setId, setX, setY, setMapname bool
	pos := types.Position{}
	item := types.Item{}

	if argc == 1 {
		lines, ok := vm.API.PartialEval(vm.API.QuoteOrElse(args.Value))
		if !ok {
			gelo.ArgumentError(vm, "item", "[properties]", args)
		}
		for line := lines; line != nil; line = line.Next {
			lineList := line.Value.(*gelo.List)
			switch vm.API.SymbolOrElse(lineList.Value).String() {
			case "id:":
				if setId {
					gelo.RuntimeError(vm, "Attempted to set id twice")
				}
				item.Id = vm.API.SymbolOrElse(lineList.Next.Value).String()
				setId = true

			case "mapname:":
				if setMapname {
					gelo.RuntimeError(vm, "Attempted to set mapname twice")
				}
				pos.Mapid = vm.API.SymbolOrElse(lineList.Next.Value).String()
				setMapname = true

			case "x:":
				if setX {
					gelo.RuntimeError(vm, "Attempted to set x twice")
				}
				num := ensureInt(vm, lineList.Next.Value)
				item.X = num
				pos.X = num
				setX = true

			case "y:":
				if setY {
					gelo.RuntimeError(vm, "Attempted to set y twice")
				}
				num := ensureInt(vm, lineList.Next.Value)
				item.Y = num
				pos.Y = num
				setY = true
			}
		}
	}

	if !setId {
		item.Id = vm.API.SymbolOrElse(vm.Ns.LookupOrElse(gelo.Convert("id"))).String()
	}
	if !setMapname {
		pos.Mapid = vm.API.SymbolOrElse(vm.Ns.LookupOrElse(gelo.Convert("mapname"))).String()
	}
	if !setX {
		num := ensureInt(vm, vm.Ns.LookupOrElse(gelo.Convert("x")))
		pos.X = num
		item.X = num
	}
	if !setY {
		num := ensureInt(vm, vm.Ns.LookupOrElse(gelo.Convert("y")))
		pos.Y = num
		item.Y = num
	}

	w.MapItems[pos] = item
	return gelo.Null
}

func (w *constArgs) ItemBlueprint(vm *gelo.VM, args *gelo.List, argc uint) gelo.Word {
	if argc > 1 {
		gelo.ArgumentError(vm, "itemb", "[properties]", args)
	}

	var setId, setName, setType, setVariation, setDescription, setEquippable, setCode bool
	var id string
	blueprint := types.ItemBlueprint{
		Effect: &itemRunner{
			world: w.world,
			vm:    vm,
			ns:    vm.Ns.Locals(vm.Ns.Depth() - 3), // ignore the namespaces with the language and the "extraapi"
		},
	}

	if argc == 1 {
		lines, ok := vm.API.PartialEval(vm.API.QuoteOrElse(args.Value))
		if !ok {
			gelo.ArgumentError(vm, "itemb", "[properties]", args)
		}
		for line := lines; line != nil; line = line.Next {
			lineList := line.Value.(*gelo.List)
			switch vm.API.SymbolOrElse(lineList.Value).String() {
			case "id:":
				if setId {
					gelo.RuntimeError(vm, "Attempted to set item id twice.")
				}
				id = vm.API.SymbolOrElse(lineList.Next.Value).String()
				setId = true

			case "name:":
				if setName {
					gelo.RuntimeError(vm, "Attempted to set item name twice.")
				}
				blueprint.Name = vm.API.SymbolOrElse(lineList.Next.Value).String()
				setName = true

			case "type:":
				if setType {
					gelo.RuntimeError(vm, "Attempted to set item type twice.")
				}
				blueprint.Type = vm.API.SymbolOrElse(lineList.Next.Value).String()
				setType = true

			case "variation:":
				if setVariation {
					gelo.RuntimeError(vm, "Attempted to set item variation twice.")
				}
				num := ensureInt(vm, lineList.Next.Value)
				blueprint.Variation = num
				setVariation = true

			case "description:":
				if setDescription {
					gelo.RuntimeError(vm, "Attempted to set item description twice.")
				}
				blueprint.Description = vm.API.SymbolOrElse(lineList.Next.Value).String()
				setDescription = true

			case "equippable:":
				if setEquippable {
					gelo.RuntimeError(vm, "Attempted to set item equippable twice.")
				}
				switch vm.API.SymbolOrElse(lineList.Next.Value).String() {
				case "none":
					blueprint.Equippable = types.EQUIP_NONE

				case "weapon":
					blueprint.Equippable = types.EQUIP_WEAP

				case "armor":
					blueprint.Equippable = types.EQUIP_ARMO

				default:
					gelo.RuntimeError(vm, "Equippable must be either \"weapon\", \"armor\" or \"none\"")
				}
				setEquippable = true

			case "code:":
				if setCode {
					gelo.RuntimeError(vm, "Attempted to set item code twice.")
				}
				blueprint.Effect.(*itemRunner).code = vm.API.QuoteOrElse(lineList.Next.Value)
				setCode = true
			}
		}
	}

	if !setId {
		id = vm.API.SymbolOrElse(vm.Ns.LookupOrElse(gelo.Convert("id"))).String()
	}
	if !setName {
		blueprint.Name = vm.API.SymbolOrElse(vm.Ns.LookupOrElse(gelo.Convert("name"))).String()
	}
	if !setType {
		blueprint.Type = vm.API.SymbolOrElse(vm.Ns.LookupOrElse(gelo.Convert("type"))).String()
	}
	if !setVariation {
		num := ensureInt(vm, vm.Ns.LookupOrElse(gelo.Convert("variation")))
		blueprint.Variation = num
	}
	if !setDescription {
		blueprint.Description = vm.API.SymbolOrElse(vm.Ns.LookupOrElse(gelo.Convert("description"))).String()
	}
	if !setEquippable {
		if symb, ok := vm.Ns.Lookup(gelo.Convert("equippable")); ok {
			switch vm.API.SymbolOrElse(symb).String() {
			case "none":
				blueprint.Equippable = types.EQUIP_NONE

			case "weapon":
				blueprint.Equippable = types.EQUIP_WEAP

			case "armor":
				blueprint.Equippable = types.EQUIP_ARMO

			default:
				gelo.RuntimeError(vm, "Equippable must be either \"weapon\" or \"armor\"")
			}
		} else {
			blueprint.Equippable = types.EQUIP_NONE
		}
	}
	if !setCode {
		blueprint.Effect.(*itemRunner).code = vm.API.QuoteOrElse(vm.Ns.LookupOrElse(gelo.Convert("code")))
	}

	w.world.ItemBlueprints[id] = blueprint
	return gelo.Null
}

func nudge(char *types.Character, val string, amount int) int {
	switch val {
	case "health":
		health := char.Health
		char.Health += amount
		if char.Health > char.MaxHealth {
			char.Health = char.MaxHealth
		}
		amount = char.Health - health

	case "maxhealth":
		char.MaxHealth += amount
		if char.Health > char.MaxHealth {
			char.Health = char.MaxHealth
		}

	case "viewdist":
		char.ViewDist += amount
	}

	return amount
}

func nudgeMessage(world *types.World, origin, target *types.Character, nudgeType string, amount int) {
	world.Updates = append(world.Updates, types.Update{
		Pos: target.Pos,
		Content: types.NudgeCharUpdate{
			Nudge:      nudgeType,
			Amount:     amount,
			OriginChar: origin.Id,
			Target:     target.Id,
		},
	})
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

	i.vm.RegisterBundle(liveBundle(i.world)) //TODO: might be interesting to optimize by adding the liveBundle after const is done and then not removing it, so it doesn't have to be readded all the time.
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
	amount := ensureInt(vm, args.Next.Value)
	actualAmount := nudge(c, val, amount)
	nudgeMessage(o.world, o.origin, c, val, actualAmount)

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
	if c.RemoveItem(bid) {
		return gelo.True
	}
	return gelo.False
}
