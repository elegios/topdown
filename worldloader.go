package main

import (
	"os"
	"path/filepath"
)

const (
	CONST_FOLDER   = "world"
	MAPS_FOLDER    = "maps"
	LIVE_FOLDER    = "live"
	INITIAL_FOLDER = "livestart"
	EXTRA_EXT      = ".ext"
	MAP_EXT        = ".png"
)

type World struct {
	itemBlueprints map[string]ItemBlueprint
	characters     map[string]*Character
	items          map[string][]Item
	props          map[string][]Prop
	maps           map[string][][]bits
	mapRoot        string
	liveRoot       string
}

func loadWorld(root string) *World {
	w := &World{
		itemBlueprints: make(map[string]ItemBlueprint),
		characters:     make(map[string]*Character),
		items:          make(map[string][]Item),
		props:          make(map[string][]Prop),
		maps:           make(map[string][][]bits),
	}
	w.mapRoot = filepath.Join(root, CONST_FOLDER, MAPS_FOLDER)
	w.liveRoot = filepath.Join(root, LIVE_FOLDER)

	liveRoot := w.liveRoot
	if _, err := os.Stat(w.liveRoot); os.IsNotExist(err) {
		w.liveRoot = filepath.Join(root, CONST_FOLDER, INITIAL_FOLDER)
	}
	filepath.Walk(w.liveRoot, w.loadLiveData)
	w.liveRoot = liveRoot

	filepath.Walk(w.mapRoot, w.loadMapData)

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
	w.maps[name] = parseMap(path)
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
	//TODO
	return nil
}

func getName(root, path string) string {
	name, _ := filepath.Rel(root, path)
	name = filepath.ToSlash(name)
	name = name[:len(name)-len(filepath.Ext(name))]
	return name
}
