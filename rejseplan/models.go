package rejseplan

import "encoding/json"

type LocationResponse struct {
	LocationList LocationList
}

type LocationList struct {
	StopLocation json.RawMessage
}

type Stop struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (s Stop) Title() string {
	return s.Name
}

type DepartureResponse struct {
	DepartureBoard DepartureBoard
}

type DepartureBoard struct {
	Departure []Departure
}

type Departure struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Stop        string `json:"stop"`
	PlannedTime string `json:"time"`
	PlannedDate string `json:"date"`
	ActualTime  string `json:"rtTime"`
	ActualDate  string `json:"rtDate"`
	ID          string `json:"id"`
	Line        string `json:"line"`
	Messages    string `json:"messages"`
	FinalStop   string `json:"finalStop"`
	Direction   string `json:"direction"`
}
