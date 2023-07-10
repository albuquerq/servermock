package handlers

import "net/http"

type StackHandler struct {
	stack []http.HandlerFunc
	base  http.HandlerFunc
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

	if len(sh.stack) == 0 {
		return sh.base
	}
	idx := len(sh.stack) - 1

	h := sh.stack[idx]
	sh.stack = sh.stack[:idx]

	return h
}

func (sh *StackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if sh == nil {
		sh = newEmptyStackHandler()
	}
	sh.next().ServeHTTP(w, r)
}

func (sh *StackHandler) Add(h http.HandlerFunc) {
	if h == nil {
		return
	}
	if sh == nil {
		sh = newEmptyStackHandler()
	}
	sh.stack = append(sh.stack, h)
}

func (sh *StackHandler) Reset() {
	if sh == nil {
		sh = newEmptyStackHandler()
	}
	sh.stack = []http.HandlerFunc{}
}
