package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

func (hc *HandlerContext) GetTournamentInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tournamentID := vars["tournament_id"]

	params := url.Values{}
	body, err := Get("https://api.rating.chgk.net/tournaments/"+tournamentID, params)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var tournamentInfo TournamentInfo
	err = json.Unmarshal(body, &tournamentInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(tournamentInfo)
}
