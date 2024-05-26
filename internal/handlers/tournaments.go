package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"
)

func GetTournaments(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	timeFormat := "2006-01-02T15:04:05-07:00"
	baseURL := "https://api.rating.chgk.net"
	countryId := queryParams.Get("countryId")

	tournamentsMap, err := fetchTournamentsData(baseURL, timeFormat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	townsIds, err := getTownsIds(baseURL, countryId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tournaments, err := filterByCountry(tournamentsMap, townsIds, baseURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	orderByDateStart(tournaments, timeFormat)
	truncateTime(tournaments)

	json.NewEncoder(w).Encode(tournaments)
}

func fetchTournamentsData(baseURL string, timeFormat string) (map[int]Tournament, error) {
	t := time.Now().UTC()
	t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	currentDate := t.AddDate(0, 0, 1).Format(timeFormat)
	upperBound := t.AddDate(0, 3, 0).Format(timeFormat)

	params := url.Values{}
	params.Add("type", "2")
	params.Add("itemsPerPage", "100")
	params.Add("dateEnd[after]", currentDate)
	params.Add("dateEnd[before]", upperBound)

	body, err := Get(baseURL+"/tournaments", params)
	if err != nil {
		return nil, err
	}
	var tournaments []Tournament
	err = json.Unmarshal(body, &tournaments)
	if err != nil {
		return nil, err
	}

	tournamentsMap := make(map[int]Tournament, len(tournaments))
	for _, tournament := range tournaments {
		tournamentsMap[tournament.ID] = tournament
	}

	return tournamentsMap, nil
}

func getTownsIds(baseURL string, countryId string) (map[int]Town, error) {
	params := url.Values{}
	params.Add("country", countryId)
	params.Add("itemsPerPage", "500")

	body, err := Get(baseURL+"/towns", params)
	if err != nil {
		return nil, err
	}
	var towns []Town
	err = json.Unmarshal(body, &towns)
	if err != nil {
		return nil, err
	}

	townsMap := make(map[int]Town, len(towns))
	for _, town := range towns {
		townsMap[town.ID] = town
	}

	return townsMap, nil
}

func filterByCountry(tournaments map[int]Tournament, townsMap map[int]Town, baseURL string) ([]Tournament, error) {
	result := make([]Tournament, 0, len(tournaments))
	ch := make(chan []byte)
	cherr := make(chan error)
	var wg sync.WaitGroup
	for _, tournament := range tournaments {
		wg.Add(1)
		go GetAsync(baseURL+"/tournaments/"+fmt.Sprint(tournament.ID), url.Values{}, ch, cherr, &wg)
	}
	go func() {
		wg.Wait()
		close(ch)
		close(cherr)
	}()

	for body := range ch {
		if body == nil {
			continue
		}
		ttData := TournamentTownData{}
		err := json.Unmarshal(body, &ttData)
		if err != nil {
			return nil, err
		}
		if town, ok := townsMap[ttData.TownId]; ok {
			tournament := tournaments[ttData.TournamentId]
			tournament.TownName = town.Name
			tournaments[ttData.TournamentId] = tournament
			result = append(result, tournaments[ttData.TournamentId])
		}
	}
	return result, nil
}

func orderByDateStart(tournaments []Tournament, timeFormat string) {
	sort.Slice(tournaments, func(i, j int) bool {
		startI, errI := time.Parse(timeFormat, tournaments[i].DateStart)
		startJ, errJ := time.Parse(timeFormat, tournaments[j].DateStart)
		if errI != nil || errJ != nil {
			return false
		}
		return startI.Before(startJ)
	})
}

func truncateTime(tournaments []Tournament) {
	for i := range tournaments {
		tournaments[i].DateStart = strings.Split(tournaments[i].DateStart, "T")[0]
		tournaments[i].DateEnd = strings.Split(tournaments[i].DateEnd, "T")[0]
	}
}
