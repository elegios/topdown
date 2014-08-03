package script

import (
	"code.google.com/p/gelo"
)

func ensureInt(vm *gelo.VM, word gelo.Word) int {
	num, ok := vm.API.NumberOrElse(word).Int()
	if !ok {
		gelo.TypeMismatch(vm, "int", word)
	}
	return int(num)
}
