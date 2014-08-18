package types

import (
	"errors"
	"os"
	"path/filepath"
)

const (
	MAIN_NAME   = "main"
	LIVE_FOLDER = "live"

	MODULE_FOLDER  = "modules"
	MAPS_FOLDER    = "maps"
	STORY_FOLDER   = "stories"
	PARTIAL_FOLDER = "partials"

	CHARACTER_FILE    = "characters"
	ANNOUNCEMENT_FILE = "announcements"
	STORIES_FILE      = "stories"
	MODULES_FILE      = "modules"
	PARTIALS_FILE     = "partials"

	CONST_SCRIPT_EXT = ".lua"
	STORY_SCRIPT_EXT = ".lua"
	MAP_EXT          = ".png"
	PARTIAL_EXT      = ".png"
	ITEM_EXT         = ".ite"
	PROP_EXT         = ".prp"
)

type VM interface {
	RunConstantScript(path, name string) error
	RunStoryScript(module, name, path string, first bool)
}

type World struct {
	saved
	constant
	Updates        []Update
	vm             VM
	root           string
	tempModuleRoot string
}

func (w *World) Load(vm VM, root string) (err error) {
	w.vm = vm
	w.root = root
	w.Updates = make([]Update, 0)
	w.initConstant()

	err = w.loadSaved(filepath.Join(root, LIVE_FOLDER))
	if err == nil {
		for _, m := range w.Modules {
			if err = w.applyModule(m, false); err != nil {
				return
			}
		}
		for _, p := range w.Partials {
			if err = w.applyPartial(p); err != nil {
				return
			}
		}
		numMods := len(w.Modules)
		numPars := len(w.Partials)
		for story := range w.Stories {
			w.runStory(story, false)
		}
		if len(w.Modules) != numMods || len(w.Partials) != numPars {
			return errors.New("A story in recovery applied a partial or module")
		}

	} else {
		if os.IsNotExist(err) {
			w.initSaved(filepath.Join(root, LIVE_FOLDER))
			return w.ApplyModule(MAIN_NAME)
		} else {
			return
		}
	}

	return
}

func (w *World) ReloadMaps() (err error) {
	w.initConstant()

	for _, m := range w.Modules {
		if err = w.applyModule(m, false); err != nil {
			return
		}
	}
	for _, p := range w.Partials {
		if err = w.applyPartial(p); err != nil {
			return
		}
	}
	return
}

func (w *World) SaveWorld() error {
	return w.save()
}

type constant struct {
	Maps           map[string][][]Bits
	ItemBlueprints map[string]ItemBlueprint
	MapTransitions map[Position]Position
}

func (c *constant) initConstant() {
	*c = constant{
		Maps:           make(map[string][][]Bits),
		ItemBlueprints: make(map[string]ItemBlueprint),
		MapTransitions: make(map[Position]Position),
	}
}

func d(err error) {
	if err != nil {
		panic(err)
	}
}
