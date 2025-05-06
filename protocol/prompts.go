package protocol

import (
	"errors"
	"maps"
	"slices"
	"sync"
)

type prompts struct {
	sync.RWMutex
	ps map[string]IPrompt
}

func (p *prompts) add(prompt IPrompt) {
	p.Lock()
	defer p.Unlock()
	promptName := prompt.getName()
	if _, ok := p.ps[promptName]; ok {
		return
	}
	p.ps[promptName] = prompt
}

func (p *prompts) fetch(name string) (IPrompt, error) {
	p.RLock()
	defer p.RUnlock()
	if prompt, ok := p.ps[name]; !ok {
		return nil, errors.New("prompt not found")
	} else {
		return prompt, nil
	}
}

func (p *prompts) remove(name string) error {
	p.Lock()
	defer p.Unlock()
	if _, ok := p.ps[name]; ok {
		delete(p.ps, name)
		return nil
	} else {
		return errors.New("prompt not found")
	}
}

func (p *prompts) reset(name string, prompt IPrompt) {
	p.Lock()
	defer p.Unlock()
	if _, ok := p.ps[name]; ok {
		p.ps[name] = prompt
	}
}

func (p *prompts) fetchAll() ([]IPrompt, error) {
	p.RLock()
	defer p.RUnlock()
	if p.ps == nil {
		return nil, errors.New("prompts is empty")
	}
	list := make([]IPrompt, 0, len(p.ps))
	for _, prompt := range p.ps {
		list = append(list, prompt)
	}
	return list, nil
}

func (p *prompts) names() []string {
	if len(p.ps) == 0 {
		return nil
	}
	p.RLock()
	defer p.RUnlock()
	return slices.Sorted(maps.Keys(p.ps))
}
