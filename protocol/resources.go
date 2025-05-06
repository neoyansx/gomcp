package protocol

import (
	"errors"
	"maps"
	"slices"
	"sync"
)

type (
	resources struct {
		sync.RWMutex
		rs map[string]IResource
	}
)

func (r *resources) add(uri string, resource IResource) {
	r.Lock()
	defer r.Unlock()
	r.rs[uri] = resource
}

func (r *resources) fetch(uri string) (IResource, error) {
	r.RLock()
	defer r.RUnlock()
	if resource, ok := r.rs[uri]; !ok {
		return nil, errors.New("resource not found")
	} else {
		return resource, nil
	}
}

func (r *resources) remove(uri string) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.rs[uri]; ok {
		delete(r.rs, uri)
		return nil
	} else {
		return errors.New("tool not found")
	}
}

func (r *resources) resetResource(uri string, resource IResource) {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.rs[uri]; ok {
		r.rs[uri] = resource
	}
}

func (r *resources) exists(uri string) bool {
	r.RLock()
	defer r.RUnlock()
	_, ok := r.rs[uri]
	return ok
}

func (r *resources) fetchAll() ([]IResource, error) {
	r.RLock()
	defer r.RUnlock()
	if r.rs == nil {
		return nil, errors.New("resources is empty")
	}
	list := make([]IResource, 0, len(r.rs))
	for _, resource := range r.rs {
		list = append(list, resource)
	}
	return list, nil
}

func (r *resources) uris() []string {
	if len(r.rs) == 0 {
		return nil
	}
	r.RLock()
	defer r.RUnlock()
	return slices.Sorted(maps.Keys(r.rs))
}
