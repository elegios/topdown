package types

import (
	"path/filepath"
)

type storyKey struct {
	module, name string
}

func (s *storyKey) String() string {
	return filepath.Join(s.module, STORY_FOLDER, s.name)
}

func (w *World) RunStory(module, name string) {
	w.runStory(storyKey{module, name}, true)
}

func (w *World) runStory(k storyKey, first bool) {
	if first {
		w.Stories[k] = make(map[string]interface{})
	}
	w.vm.RunStoryScript(k.module, k.name, filepath.Join(w.root, MODULE_FOLDER, k.String()+STORY_SCRIPT_EXT), first)
}

func (w *World) RemoveStory(module, name string) {
	delete(w.Stories, storyKey{module, name})
}

func (w *World) ReadStoryValue(module, name, key string) interface{} {
	return w.Stories[storyKey{module, name}][key]
}

func (w *World) SetStoryValue(module, name, key string, value interface{}) {
	w.Stories[storyKey{module, name}][key] = value
}
