package protocol

import (
	"errors"
	"sync"
)

type tools struct {
	sync.RWMutex
	ts map[string]ITool
}

func (t *tools) add(tool ITool) {
	t.Lock()
	defer t.Unlock()
	toolName := tool.getName()
	if _, ok := t.ts[toolName]; ok {
		return
	}
	t.ts[toolName] = tool
}

func (t *tools) addMore(tools ...ITool) {
	if len(tools) == 0 {
		return
	}
	t.Lock()
	defer t.Unlock()
	for _, tool := range tools {
		if _, ok := t.ts[tool.getName()]; ok {
			return
		}
		t.ts[tool.getName()] = tool
	}
}

func (t *tools) fetch(name string) (ITool, error) {
	t.RLock()
	defer t.RUnlock()
	if tool, ok := t.ts[name]; !ok {
		return nil, errors.New("tool not found")
	} else {
		return tool, nil
	}
}

func (t *tools) remove(name string) error {
	t.Lock()
	defer t.Unlock()
	if _, ok := t.ts[name]; ok {
		delete(t.ts, name)
		return nil
	} else {
		return errors.New("tool not found")
	}
}

func (t *tools) reset(name string, tool ITool) {
	t.Lock()
	defer t.Unlock()
	if _, ok := t.ts[name]; ok {
		t.ts[name] = tool
	}
}

func (t *tools) fetchAll() ([]ITool, error) {
	t.RLock()
	defer t.RUnlock()
	if t.ts == nil {
		return nil, errors.New("tools is empty")
	}
	list := make([]ITool, 0, len(t.ts))
	for _, tool := range t.ts {
		list = append(list, tool)
	}
	return list, nil
}
