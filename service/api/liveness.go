package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// liveness is an HTTP handler that checks the API server status. If the server cannot serve requests (e.g., some
// resources are not ready), this should reply with HTTP Status 500. Otherwise, with HTTP Status 200
func (rt *_router) liveness(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("WE GUCCI"))
}
