package types

import (
	"os"
	"path/filepath"
)

func (w *World) ApplyModule(module string) error {
	return w.applyModule(module, true)
}

func (w *World) applyModule(module string, first bool) (err error) {
	w.tempModuleRoot = filepath.Join(w.root, MODULE_FOLDER, module)

	err = filepath.Walk(filepath.Join(w.tempModuleRoot, MAPS_FOLDER), w.loadMapData)
	if err != nil {
		return
	}

	if first {
		err = w.RunStory(module, "main")
		if err != nil {
			return
		}
	}

	w.Modules = append(w.Modules, module)
	return
}

func (w *World) loadMapData(path string, info os.FileInfo, _ error) error {
	if info.IsDir() {
		return nil
	}

	name := getName(filepath.Join(w.tempModuleRoot, MAPS_FOLDER), path)

	switch filepath.Ext(path) {
	case MAP_EXT:
		return w.loadMap(path, name)

	case CONST_SCRIPT_EXT:
		return w.vm.RunConstantScript(path, name)

	default:
		return nil
	}
}
