package handlers

import "net/http"

type StackHandler struct {
	pool []http.HandlerFunc
	base http.HandlerFunc
}

func NewStackHandler(base http.HandlerFunc) *StackHandler {
	if base == nil {
		return newEmptyStackHandler()
	}
	return &StackHandler{
		base: base,
	}
}

func newEmptyStackHandler() *StackHandler {
	return &StackHandler{
		base: func(http.ResponseWriter, *http.Request) {},
	}
}

func (sh *StackHandler) next() http.HandlerFunc {
	if sh == nil {
		sh = newEmptyStackHandler()
	}

	if len(sh.pool) == 0 {
		return sh.base
	}
	idx := len(sh.pool) - 1

	h := sh.pool[idx]
	sh.pool = sh.pool[:idx]

	return h
}

func (hp *StackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if hp == nil {
		hp = newEmptyStackHandler()
	}
	hp.next().ServeHTTP(w, r)
}

func (hp *StackHandler) Add(h http.HandlerFunc) {
	if h == nil {
		return
	}
	if hp == nil {
		hp = newEmptyStackHandler()
	}
	hp.pool = append(hp.pool, h)
}

func (hp *StackHandler) Reset() {
	if hp == nil {
		hp = newEmptyStackHandler()
	}
	hp.pool = []http.HandlerFunc{}
}
