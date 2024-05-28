package handlers

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

func (hc *HandlerContext) GetTeams(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tournamentID := vars["tournament_id"]

	params := url.Values{}
	body, err := Get("https://api.rating.chgk.net/tournaments/"+tournamentID+"/results", params)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var results []TeamData
	err = json.Unmarshal(body, &results)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	teams := make([]map[string]interface{}, len(results))
	for i, result := range results {
		teams[i] = map[string]interface{}{
			"questionsTotal": result.QuestionsTotal,
			"position":       result.Position,
			"picked":         false,
			"name":           result.Current.Name,
			"cost":           1 + rand.Intn(100),
		}
	}

	json.NewEncoder(w).Encode(teams)
}
