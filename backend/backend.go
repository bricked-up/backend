// Package backend provides the backend infrastructure (route handling + database)
// for the Bricked-Up website.
package backend

import (
	"brickedup/backend/endpoints"
	"database/sql"
	"net/http"
)

// MainHandler checks if the request URL matches a known endpoint.
// If it does, the corresponding handler is called; otherwise, it returns a 404 error.
func MainHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if handler, ok := endpoints.Endpoints[r.URL.Path]; ok {
		handler(db, w, r)
		return
	}
	http.NotFound(w, r)
}
