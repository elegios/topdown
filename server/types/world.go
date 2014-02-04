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

type World struct {
	ItemBlueprints map[string]ItemBlueprint
	Charids        map[string]*Character
	MapCharacters  map[Position]*Character
	MapItems       map[Position]Item
	MapProps       map[Position]Prop
	Maps           map[string][][]Bits
	mapRoot        string
	liveRoot       string
}

func LoadWorld(root string) *World {
	w := &World{
		ItemBlueprints: make(map[string]ItemBlueprint),
		Charids:        make(map[string]*Character),
		MapCharacters:  make(map[Position]*Character),
		MapItems:       make(map[Position]Item),
		MapProps:       make(map[Position]Prop),
		Maps:           make(map[string][][]Bits),

		mapRoot:  filepath.Join(root, CONST_FOLDER, MAPS_FOLDER),
		liveRoot: filepath.Join(root, LIVE_FOLDER, MAPS_FOLDER),
	}

	filepath.Walk(w.mapRoot, w.loadMapData)

	liveRoot := w.liveRoot
	if _, err := os.Stat(w.liveRoot); os.IsNotExist(err) {
		w.liveRoot = filepath.Join(root, CONST_FOLDER, INITIAL_FOLDER)
	}
	filepath.Walk(w.liveRoot, w.loadLiveData)
	w.liveRoot = liveRoot

	w.loadCharacters(filepath.Join(root, LIVE_FOLDER, CHARACTER_FILE))

	return w
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

func (w *World) loadMap(path, name string) error {
	w.Maps[name] = parseMap(path)
	return nil
}

func (w *World) loadMapExtra(path, name string) error {
	//TODO
	return nil
}

func (w *World) loadLiveData(path string, info os.FileInfo, err error) error {
	if info.IsDir() {
		return nil
	}

	name := getName(w.liveRoot, path)

	switch filepath.Ext(path) {
	case ITEM_EXT:
		return w.loadItems(path, name)

	case PROP_EXT:
		return w.loadProps(path, name)

	default:
		return nil
	}
}

func getName(root, path string) string {
	name, _ := filepath.Rel(root, path)
	name = filepath.ToSlash(name)
	name = name[:len(name)-len(filepath.Ext(name))]
	return name
}
