package routes

import (
	"net/http"

	"github.com/akim-malyshchyk/fantasy-backend/internal/handlers"
	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func NewRouter() http.Handler {
	r := mux.NewRouter()
	r.Use(responseContentTypeMiddleware)

	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/tournaments/{tournament_id:[0-9]+}/teams", handlers.GetTeams).Methods("GET")

	corsOrigins := gorillaHandlers.AllowedOrigins([]string{"*"})
	corsMethods := gorillaHandlers.AllowedMethods([]string{"GET"})
	corsHeaders := gorillaHandlers.AllowedHeaders([]string{"Content-Type"})

	return gorillaHandlers.CORS(corsOrigins, corsMethods, corsHeaders)(r)
}
