package script

import (
	"code.google.com/p/gelo"
	"code.google.com/p/gelo/extensions"
	"github.com/elegios/topdown/server/types"
)

func (o *onAttackVals) Speak(vm *gelo.VM, args *gelo.List, argc uint) gelo.Word {
	speech := vm.API.SymbolOrElse(args.Value).String()
	update := types.Update{
		Pos: o.origin.Pos,
		Content: types.SpeechCharUpdate{
			Speech:    speech,
			Character: o.origin.Id,
		},
	}
	o.world.Updates = append(o.world.Updates, update)
	return gelo.Null
}

var (
	speakAtParser = extensions.MakeOrElseArgParser("speech 'at name x y")
)

func (w *vmworld) SpeakAt(vm *gelo.VM, args *gelo.List, argc uint) gelo.Word {
	dict := speakAtParser(vm, args)
	speech := vm.API.SymbolOrElse(dict["speech"]).String()
	name := vm.API.SymbolOrElse(dict["name"]).String()
	x := ensureInt(vm, dict["x"])
	y := ensureInt(vm, dict["y"])

	pos := types.Position{
		Mapid: name,
		X:     x,
		Y:     y,
	}
	update := types.Update{
		Pos: pos,
		Content: types.SpeechPointUpdate{
			Speech: speech,
			Pos:    pos,
		},
	}
	w.Updates = append(w.Updates, update)

	return gelo.Null
}
