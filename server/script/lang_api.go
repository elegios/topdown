package script

import (
	"code.google.com/p/gelo"
	"code.google.com/p/gelo/commands"
	"code.google.com/p/gelo/extensions"
)

var (
	langCommands = []map[string]interface{}{
		commands.ControlCommands,
		map[string]interface{}{
			"foreach": ForEach,
			"list":    List,
		},
	}
)

var (
	forEachParser     = extensions.MakeOrElseArgParser("[[name '<- ary]|arys] 'in block")
	forEachLineParser = extensions.MakeOrElseArgParser("name '<- ary")
)

func ForEach(vm *gelo.VM, args *gelo.List, argc uint) gelo.Word {
	dict := forEachParser(vm, args)
	lists := make(map[string]*gelo.List)
	if arys, ok := dict["arys"]; ok {
		lines, ok := vm.API.PartialEval(vm.API.QuoteOrElse(arys))
		if !ok {
			gelo.ArgumentError(vm, "foreach", "{name <- ary; name <- ary...} in quote", args)
		}
		for line := lines; line != nil; line = line.Next {
			linestuff := forEachLineParser(vm, line.Value.(*gelo.List))
			name := vm.API.SymbolOrElse(linestuff["name"]).String()
			lists[name] = vm.API.ListOrElse(linestuff["ary"])
		}
	} else {
		name := vm.API.SymbolOrElse(dict["name"]).String()
		lists[name] = vm.API.ListOrElse(dict["ary"])
	}

	for cont := true; cont; {
		cont = false
		for name, list := range lists {
			vm.Ns.Set(0, gelo.Convert(name), list.Value)
			if list.Next != nil {
				cont = true
				lists[name] = list.Next
			}
		}

		vm.API.InvokeCmdOrElse(dict["block"], gelo.EmptyList)
	}
	return gelo.Null
}

func List(_ *gelo.VM, args *gelo.List, _ uint) gelo.Word {
	return args
}
