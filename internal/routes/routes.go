package routes

import (
	"database/sql"
	"net/http"

	"github.com/akim-malyshchyk/fantasy-backend/internal/handlers"
	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func NewRouter(db *sql.DB) http.Handler {
	r := mux.NewRouter()
	r.Use(responseContentTypeMiddleware)

	hc := &handlers.HandlerContext{DB: db, BaseUrl: "https://api.rating.chgk.net"}

	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/tournaments/{tournament_id:[0-9]+}", hc.GetTournamentInfo).Methods("GET")
	api.HandleFunc("/tournaments/{tournament_id:[0-9]+}/teams", hc.GetTeams).Methods("GET")
	api.HandleFunc("/tournaments", hc.GetTournaments).Methods("GET")

	corsOrigins := gorillaHandlers.AllowedOrigins([]string{"*"})
	corsMethods := gorillaHandlers.AllowedMethods([]string{"GET"})
	corsHeaders := gorillaHandlers.AllowedHeaders([]string{"Content-Type"})

	return gorillaHandlers.CORS(corsOrigins, corsMethods, corsHeaders)(r)
}
