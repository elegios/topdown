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

	CONST_SCRIPT_EXT = ".lua"
	LIVE_SCRIPT_EXT  = ".lua"
	MAP_EXT          = ".png"
	ITEM_EXT         = ".ite"
	PROP_EXT         = ".prp"
)

type VM interface {
	RunConstantScript(path, name string) error
	RunLiveScript(path string) error
}

type World struct {
	ItemBlueprints map[string]ItemBlueprint
	Charids        map[string]*Character
	MapCharacters  map[Position]*Character
	MapItems       map[Position]Item
	MapProps       map[Position]Prop
	MapTransitions map[Position]Position
	Updates        []Update
	Maps           map[string][][]Bits
	mapRoot        string
	liveMapRoot    string
	liveRoot       string
	vm             VM
}

func LoadWorld(w *World, vm VM, root string) *World {
	*w = World{
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
	w.MapTransitions = make(map[Position]Position)

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

	case CONST_SCRIPT_EXT:
		return w.loadMapScript(path, name)

	default:
		return nil
	}
}

func (w *World) loadMapScript(path, name string) error {
	return w.vm.RunConstantScript(path, name)
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

	case LIVE_SCRIPT_EXT:
		return w.vm.RunLiveScript(path)

	default:
		return nil
	}
}

func (w *World) SaveWorld() {
	w.saveCharacters(filepath.Join(w.liveRoot, CHARACTER_FILE))
	w.saveItems(w.liveMapRoot)
	w.saveProps(w.liveMapRoot)
}
