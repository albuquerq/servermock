package handlers

import "net/http"

type QueueHandler struct {
	queue []http.HandlerFunc
	base  http.HandlerFunc
}

func NewQueueHandler(base http.HandlerFunc) *QueueHandler {
	if base == nil {
		return newEmptyQueueHandler()
	}
	return &QueueHandler{
		base: base,
	}
}

func newEmptyQueueHandler() *QueueHandler {
	return &QueueHandler{
		base: func(http.ResponseWriter, *http.Request) {},
	}
}

func (qh *QueueHandler) next() http.HandlerFunc {
	if qh == nil {
		qh = newEmptyQueueHandler()
	}

	if len(qh.queue) == 0 {
		return qh.base
	}

	h := qh.queue[0]
	if len(qh.queue) > 0 {
		qh.queue = qh.queue[1:]
	}
	return h
}

func (qh *QueueHandler) Add(h http.HandlerFunc) {
	if h == nil {
		return
	}
	if qh == nil {
		qh = newEmptyQueueHandler()
	}
	qh.queue = append(qh.queue, h)
}

func (qh *QueueHandler) Reset() {
	if qh == nil {
		qh = newEmptyQueueHandler()
	}
	qh.queue = []http.HandlerFunc{}
}

func (qh *QueueHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if qh == nil {
		qh = newEmptyQueueHandler()
	}
	qh.next().ServeHTTP(w, r)
}
