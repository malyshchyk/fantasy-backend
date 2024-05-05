package handlers

import (
	"encoding/json"
	"io"
	"math/rand"
	"net/http"

	"github.com/gorilla/mux"
)

type Team struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Town Town   `json:"town"`
}

type Town struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Result struct {
	Team    Team `json:"team"`
	Current Team `json:"current"`
}

func GetTeams(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tournamentID := vars["tournament_id"]

	resp, err := http.Get("https://api.rating.chgk.net/tournaments/" + tournamentID + "/results")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var results []Result
	err = json.Unmarshal(body, &results)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	teams := make([]map[string]interface{}, len(results))
	for i, result := range results {
		teams[i] = map[string]interface{}{
			"picked": false,
			"name":   result.Current.Name,
			"cost":   1 + rand.Intn(100),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teams)
}
