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
	body, err := Get(hc.BaseUrl+"/tournaments/"+tournamentID+"/results", params)

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
			"id":             result.Team.ID,
			"name":           result.Current.Name,
			"questionsTotal": result.QuestionsTotal,
			"position":       result.Position,
			"picked":         false,
			"cost":           1 + rand.Intn(100),
		}
	}

	json.NewEncoder(w).Encode(teams)
}
