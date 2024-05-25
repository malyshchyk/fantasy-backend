package handlers

type Tournament struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	DateStart string `json:"dateStart"`
	DateEnd   string `json:"dateEnd"`
}

type TournamentData struct {
	ID     int `json:"id"`
	Idtown int `json:"idtown"`
}

type Team struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Town Town   `json:"town"`
}

type Country struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Town struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Country Country `json:"country"`
}

type TeamData struct {
	Team           Team    `json:"team"`
	Current        Team    `json:"current"`
	QuestionsTotal int     `json:"questionsTotal"`
	Position       float32 `json:"position"`
}
