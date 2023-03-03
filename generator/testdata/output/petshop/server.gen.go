// Code generated at by github.com/albuquerq/servermock DO NOT EDIT.
package petshop

import (
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi/v5"
)

type FakeServer struct {
	*httptest.Server
	// handlers.
	AddPet http.HandlerFunc // [POST] /pet
}

func NewFakeServer() *FakeServer {
	srv := &FakeServer{
		// handlers.
		AddPet: AddPetStatusOk,
	}
	srv.initialize()
	return srv
}

func (fs *FakeServer) initialize() {
	router := chi.NewRouter()

	router.Method("POST", "/pet", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.AddPet(w, r)
		fs.AddPet = AddPetStatusOk
	}))

	fs.Server = httptest.NewServer(router)
}

func sendResponse(w http.ResponseWriter, payload string, statusCode int) { //nolint
	w.WriteHeader(statusCode)
	_, _ = w.Write([]byte(payload))
}
