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

func (hc *HandlerContext) GetTournaments(w http.ResponseWriter, r *http.Request) {
	timeFormat := "2006-01-02T15:04:05-07:00"
	countryIds := r.URL.Query()["countryId"]

	tournamentsMap, err := hc.fetchTournamentsData(timeFormat, countryIds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	townsIds, err := hc.getTownsIds(countryIds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tournaments, err := hc.appendTowns(tournamentsMap, townsIds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	orderByDateStart(tournaments, timeFormat)
	truncateTime(tournaments)

	json.NewEncoder(w).Encode(tournaments)
}

func (hc *HandlerContext) fetchTournamentsData(timeFormat string, countryIds []string) (map[int]Tournament, error) {
	t := time.Now().UTC()
	t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	currentDate := t.AddDate(0, 0, 1).Format(timeFormat)
	upperBound := t.AddDate(0, 3, 0).Format(timeFormat)

	params := url.Values{}
	params.Add("type", "2")
	params.Add("itemsPerPage", "100")
	params.Add("dateEnd[after]", currentDate)
	params.Add("dateEnd[before]", upperBound)
	for _, countryId := range countryIds {
		params.Add("town.country[]", countryId)
	}

	body, err := Get(hc.BaseUrl+"/tournaments", params)
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

func (hc *HandlerContext) getTownsIds(countryIds []string) (map[int]Town, error) {
	countryIdsArray := "{" + strings.Join(countryIds, ",") + "}"
	stmt, err := hc.DB.Prepare("SELECT id, name FROM town WHERE country_id = ANY($1::int[])")
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(countryIdsArray)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	towns := make(map[int]Town)

	for rows.Next() {
		var town Town
		err = rows.Scan(&town.ID, &town.Name)
		if err != nil {
			return nil, err
		}
		towns[town.ID] = town
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return towns, nil
}

func (hc *HandlerContext) appendTowns(tournamentsMap map[int]Tournament, townsMap map[int]Town) ([]Tournament, error) {
	tournaments := make([]Tournament, 0, len(tournamentsMap))
	ch := make(chan []byte)
	cherr := make(chan error)
	var wg sync.WaitGroup
	for _, tournament := range tournamentsMap {
		wg.Add(1)
		go GetAsync(hc.BaseUrl+"/tournaments/"+fmt.Sprint(tournament.ID), url.Values{}, ch, cherr, &wg)
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
			tournament := tournamentsMap[ttData.TournamentId]
			tournament.TownName = town.Name
			tournamentsMap[ttData.TournamentId] = tournament
			tournaments = append(tournaments, tournamentsMap[ttData.TournamentId])
		}
	}
	return tournaments, nil
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
