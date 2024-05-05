package routes

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter() http.Handler {
	r := mux.NewRouter()
	r.Use(responseContentTypeMiddleware)

	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/", indexHandler)

	return r
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello World!")
}
