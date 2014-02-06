package types

import (
	"os"
	"path/filepath"
)

const (
	CONST_FOLDER   = "world"
	MAPS_FOLDER    = "maps"
	LIVE_FOLDER    = "live"
	INITIAL_FOLDER = "livestart"

	CHARACTER_FILE = "characters"

	EXTRA_EXT = ".ext"
	MAP_EXT   = ".png"
	ITEM_EXT  = ".ite"
	PROP_EXT  = ".prp"
)

type VM interface {
	RunConstantScript(path string, world *World) error
}

type World struct {
	ItemBlueprints map[string]ItemBlueprint
	Charids        map[string]*Character
	MapCharacters  map[Position]*Character
	MapItems       map[Position]Item
	MapProps       map[Position]Prop
	Maps           map[string][][]Bits
	mapRoot        string
	liveMapRoot    string
	liveRoot       string
	vm             VM
}

func LoadWorld(vm VM, root string) *World {
	w := &World{
		Charids:       make(map[string]*Character),
		MapCharacters: make(map[Position]*Character),
		MapItems:      make(map[Position]Item),
		MapProps:      make(map[Position]Prop),

		mapRoot:     filepath.Join(root, CONST_FOLDER, MAPS_FOLDER),
		liveMapRoot: filepath.Join(root, LIVE_FOLDER, MAPS_FOLDER),
		liveRoot:    filepath.Join(root, LIVE_FOLDER),

		vm: vm,
	}

	d(w.LoadConstantWorld())

	liveMapRoot := w.liveMapRoot
	if _, err := os.Stat(w.liveRoot); os.IsNotExist(err) {
		w.liveMapRoot = filepath.Join(root, CONST_FOLDER, INITIAL_FOLDER, MAPS_FOLDER)
	}
	filepath.Walk(w.liveMapRoot, w.loadLiveData)
	w.liveMapRoot = liveMapRoot

	w.loadCharacters(filepath.Join(root, LIVE_FOLDER, CHARACTER_FILE))

	return w
}

func (w *World) LoadConstantWorld() error {
	w.ItemBlueprints = make(map[string]ItemBlueprint)
	w.Maps = make(map[string][][]Bits)

	return filepath.Walk(w.mapRoot, w.loadMapData)
}

func (w *World) loadMapData(path string, info os.FileInfo, err error) error {
	if info.IsDir() {
		return nil
	}

	name := getName(w.mapRoot, path)

	switch filepath.Ext(path) {
	case MAP_EXT:
		return w.loadMap(path, name)

	case EXTRA_EXT:
		return w.loadMapExtra(path, name)

	default:
		return nil
	}
}

func (w *World) loadMapExtra(path, name string) error {
	//TODO
	return w.vm.RunConstantScript(path, w)
}

func (w *World) loadLiveData(path string, info os.FileInfo, err error) error {
	if info.IsDir() {
		return nil
	}

	name := getName(w.liveMapRoot, path)

	switch filepath.Ext(path) {
	case ITEM_EXT:
		return w.loadItems(path, name)

	case PROP_EXT:
		return w.loadProps(path, name)

	default:
		return nil
	}
}

func (w *World) SaveWorld() {
	w.saveCharacters(filepath.Join(w.liveRoot, CHARACTER_FILE))
	w.saveItems(w.liveMapRoot)
	w.saveProps(w.liveMapRoot)
}
