package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

func GetTournaments(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	t := time.Now().UTC()
	timeFormat := "2006-01-02T15:04:05-07:00"
	baseURL := "https://api.rating.chgk.net"

	currentDate := t.Format(timeFormat)
	lowerBound := t.AddDate(0, -6, 0).Format(timeFormat)
	upperBound := t.AddDate(0, 6, 0).Format(timeFormat)
	countryId := queryParams.Get("countryId")
	beforeToday := queryParams.Get("beforeToday") == "true"

	params := url.Values{}
	params.Add("type", "2")
	if beforeToday {
		params.Add("dateEnd[after]", lowerBound)
		params.Add("dateStart[strictly_before]", currentDate)
	} else {
		params.Add("dateStart[after]", currentDate)
		params.Add("dateEnd[before]", upperBound)
	}

	body := Get(baseURL+"/tournaments", params, w)

	var tournaments []Tournament
	err := json.Unmarshal(body, &tournaments)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	townsIds := getTownsIds(baseURL, countryId, w)
	tournaments = filterByCountry(tournaments, townsIds, baseURL, w)

	orderByDateStart(tournaments, timeFormat, w)

	for i := range tournaments {
		tournaments[i].DateStart = strings.Split(tournaments[i].DateStart, "T")[0]
		tournaments[i].DateEnd = strings.Split(tournaments[i].DateEnd, "T")[0]
	}

	json.NewEncoder(w).Encode(tournaments)
}

func getTownsIds(baseURL string, countryId string, w http.ResponseWriter) map[int]bool {
	params := url.Values{}
	params.Add("country", countryId)
	params.Add("itemsPerPage", "500")
	body := Get(baseURL+"/towns", params, w)

	var towns []Town
	err := json.Unmarshal(body, &towns)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}

	ids := make(map[int]bool, len(towns))
	for _, town := range towns {
		ids[town.ID] = true
	}

	return ids
}

func filterByCountry(tournaments []Tournament, townsIds map[int]bool, baseURL string, w http.ResponseWriter) []Tournament {
	result := make([]Tournament, 0, len(tournaments))
	for i := range tournaments {
		body := Get(baseURL+"/tournaments/"+fmt.Sprint(tournaments[i].ID), url.Values{}, w)
		tData := TournamentData{}
		err := json.Unmarshal(body, &tData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return nil
		}
		if _, ok := townsIds[tData.Idtown]; ok {
			result = append(result, tournaments[i])
		}
	}
	return result
}

func orderByDateStart(tournaments []Tournament, timeFormat string, w http.ResponseWriter) {
	sort.Slice(tournaments, func(i, j int) bool {
		startI, errI := time.Parse(timeFormat, tournaments[i].DateStart)
		startJ, errJ := time.Parse(timeFormat, tournaments[j].DateStart)
		if errI != nil || errJ != nil {
			http.Error(w, "Error parsing tournament start dates", http.StatusInternalServerError)
			return false
		}
		return startI.Before(startJ)
	})
}
