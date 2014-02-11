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
			"let":     Let,
			"foreach": ForEach,
			"list":    List,
		},
	}
)

var (
	letParser     = extensions.MakeOrElseArgParser("[[name '= value]|values] 'in block")
	letLineParser = extensions.MakeOrElseArgParser("name '= value")
)

func Let(vm *gelo.VM, args *gelo.List, argc uint) gelo.Word {
	dict := letParser(vm, args)
	vm.Ns.Fork(nil)
	defer vm.Ns.Unfork()

	if values, ok := dict["values"]; ok {
		lines, ok := vm.API.PartialEval(vm.API.QuoteOrElse(values))
		if !ok {
			gelo.ArgumentError(vm, "let", "{name = value; name = value;...} in invokable", args)
		}
		for line := lines; line != nil; line = line.Next {
			linestuff := letLineParser(vm, line.Value.(*gelo.List))
			vm.Ns.Set(0, linestuff["name"], linestuff["value"])
		}

	} else {
		vm.Ns.Set(0, dict["name"], dict["value"])
	}

	return vm.API.InvokeCmdOrElse(dict["block"], gelo.EmptyList)
}

var (
	forEachParser     = extensions.MakeOrElseArgParser("[[name '<- ary]|arys] 'in block")
	forEachLineParser = extensions.MakeOrElseArgParser("name '<- ary")
)

func ForEach(vm *gelo.VM, args *gelo.List, argc uint) gelo.Word {
	dict := forEachParser(vm, args)
	vm.Ns.Fork(nil)
	defer vm.Ns.Unfork()
	lists := make(map[string]*gelo.List)
	if arys, ok := dict["arys"]; ok {
		lines, ok := vm.API.PartialEval(vm.API.QuoteOrElse(arys))
		if !ok {
			gelo.ArgumentError(vm, "foreach", "{name <- ary; name <- ary;...} in invokable", args)
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
