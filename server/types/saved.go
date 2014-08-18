package types

import (
	"os"
	"path/filepath"
)

type saved struct {
	Charids       map[string]*Character
	MapCharacters map[Position]*Character
	MapItems      map[Position]Item
	MapProps      map[Position]Prop
	Announcements []Announcement
	Modules       []string
	Partials      []Partial
	Stories       map[storyKey]map[string]interface{}
	liveDir       string
}
type Partial struct {
	Pos  Position
	Path string
}

func (s *saved) save() (err error) {
	err = s.saveCharacters(filepath.Join(s.liveDir, CHARACTER_FILE))
	if err != nil {
		return
	}
	err = s.saveItems(s.liveDir)
	if err != nil {
		return
	}
	err = s.saveProps(s.liveDir)
	if err != nil {
		return
	}
	err = enc(filepath.Join(s.liveDir, ANNOUNCEMENT_FILE), s.Announcements)
	if err != nil {
		return
	}
	err = enc(filepath.Join(s.liveDir, MODULES_FILE), s.Modules)
	if err != nil {
		return
	}
	err = enc(filepath.Join(s.liveDir, PARTIALS_FILE), s.Partials)
	if err != nil {
		return
	}
	return enc(filepath.Join(s.liveDir, STORIES_FILE), s.Stories)
}

func (s *saved) loadSaved(liveDir string) (err error) {
	s.liveDir = liveDir
	err = s.loadCharacters(filepath.Join(s.liveDir, CHARACTER_FILE))
	if err != nil {
		return
	}
	s.MapItems = make(map[Position]Item)
	s.MapProps = make(map[Position]Prop)
	err = filepath.Walk(s.liveDir, s.loadLiveData)
	if err != nil {
		return
	}
	err = dec(filepath.Join(s.liveDir, ANNOUNCEMENT_FILE), &s.Announcements)
	if err != nil {
		return
	}
	err = dec(filepath.Join(s.liveDir, MODULES_FILE), &s.Modules)
	if err != nil {
		return
	}
	err = dec(filepath.Join(s.liveDir, PARTIALS_FILE), &s.Partials)
	if err != nil {
		return
	}
	return dec(filepath.Join(s.liveDir, STORIES_FILE), &s.Stories)
}

func (s *saved) loadLiveData(path string, info os.FileInfo, _ error) error {
	if info.IsDir() {
		return nil
	}

	name := getName(s.liveDir, path)

	switch filepath.Ext(path) {
	case ITEM_EXT:
		return s.loadItems(path, name)

	case PROP_EXT:
		return s.loadProps(path, name)

	default:
		return nil
	}
}

func (s *saved) initSaved(liveDir string) {
	*s = saved{
		Charids:       make(map[string]*Character),
		MapCharacters: make(map[Position]*Character),
		MapItems:      make(map[Position]Item),
		MapProps:      make(map[Position]Prop),
		Announcements: make([]Announcement, 0),
		Modules:       make([]string, 0),
		Partials:      make([]Partial, 0),
		Stories:       make(map[storyKey]map[string]interface{}),
		liveDir:       liveDir,
	}
}
