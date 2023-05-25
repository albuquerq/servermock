package temp

import "net/http"

type StackHandler struct {
	pool []http.HandlerFunc
	base http.HandlerFunc
}

func NewStackHandler(base http.HandlerFunc) *StackHandler {
	if base == nil {
		base = func(http.ResponseWriter, *http.Request) {}
	}
	return &StackHandler{
		base: base,
	}
}

func (hp *StackHandler) pop() http.HandlerFunc {
	if hp == nil {
		hp = newEmpty()
	}

	if len(hp.pool) == 0 {
		return hp.base
	}
	idx := len(hp.pool) - 1

	h := hp.pool[idx]
	hp.pool = hp.pool[:idx]

	return h
}

func (hp *StackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if hp == nil {
		hp = newEmpty()
	}
	hp.pop().ServeHTTP(w, r)
}

func (hp *StackHandler) Push(h http.HandlerFunc) {
	if h == nil {
		return
	}
	if hp == nil {
		hp = newEmpty()
	}
	hp.pool = append(hp.pool, h)
}

func (hp *StackHandler) Reset() {
	if hp == nil {
		hp = newEmpty()
	}
	hp.pool = []http.HandlerFunc{}
}

func newEmpty() *StackHandler {
	return &StackHandler{
		base: func(http.ResponseWriter, *http.Request) {},
	}
}
