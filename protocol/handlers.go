package protocol

import "sync"

type (
	handler func(req Request) *Response

	handlers struct {
		sync.RWMutex
		hs map[string]handler
	}
)

func (a *handlers) add(method string, hd handler) {
	a.Lock()
	defer a.Unlock()
	_, ok := a.hs[method]
	if ok {
		return
	}
	a.hs[method] = hd
}

func (a *handlers) set(method string, hd handler) {
	a.Lock()
	defer a.Unlock()
	a.hs[method] = hd
}

func (a *handlers) get(method string) handler {
	a.RLock()
	defer a.RUnlock()
	hd, ok := a.hs[method]
	if !ok {
		return nil
	}
	return hd
}

func (a *handlers) invoke(req Request) *Response {
	hd := a.get(req.Method)
	if hd == nil {
		return nil
	}
	return hd(req)
}

func (a *handlers) remove(method string) {
	a.Lock()
	defer a.Unlock()
	_, ok := a.hs[method]
	if !ok {
		return
	}
	delete(a.hs, method)
}
