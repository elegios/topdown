package types

import (
	"path/filepath"
)

func (w *World) RunStory(module, name string) error {
	return w.runStory(filepath.Join(module, STORY_FOLDER, name), true)
}

func (w *World) runStory(relPath string, first bool) (err error) {
	if first {
		w.Stories[relPath] = make(map[string]interface{})
	}
	return w.vm.RunStoryScript(filepath.Join(w.root, MODULE_FOLDER, relPath+STORY_SCRIPT_EXT), first)
}
