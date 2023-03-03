// Code generated at by github.com/albuquerq/servermock DO NOT EDIT.
package petshop

import (
	"net/http"

	"github.com/albuquerq/fakeserver/petshop/petshpdata"
)

// AddPetStatusOk Description here.
func AddPetStatusOk(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Header().Set("x-request-id", "1014a3f5-c703-44b7-8752-dfcefe497f68")
	sendResponse(w, petshpdata.AddPetStatusOk, 200)
}

func AddPetStatusInvalidInput(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("content-type", "application/json")
	sendResponse(w, petshpdata.AddPetStatusInvalidInput, 405)
}
