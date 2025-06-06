package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// liveness is an HTTP handler that checks the API server status.
// If the server cannot serve requests (e.g., some resources are not ready), respond with HTTP 500.
// Otherwise, respond with HTTP 200.
func (rt *_router) liveness(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Future: Uncomment and use this if you want DB liveness checking
	// if err := rt.DB.Ping(); err != nil {
	//     w.WriteHeader(http.StatusInternalServerError)
	//     return
	// }

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
