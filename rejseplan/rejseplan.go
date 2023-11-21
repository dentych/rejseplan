package rejseplan

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

var baseUrl = "http://xmlopen.rejseplanen.dk/bin/rest.exe"

func GetStops(input string) []Stop {
	resp, err := http.DefaultClient.Get(baseUrl + "/location?format=json&input=" + url.QueryEscape(input))
	if err != nil {
		panic(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var response LocationResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		panic(err)
	}

	var output []Stop
	if bytes.HasPrefix(response.LocationList.StopLocation, []byte("{")) {
		var stop Stop
		err = json.Unmarshal(response.LocationList.StopLocation, &stop)
		if err != nil {
			panic(err)
		}
		output = []Stop{stop}
	} else {
		err = json.Unmarshal(response.LocationList.StopLocation, &output)
		if err != nil {
			panic(err)
		}
	}

	return output
}

func GetDepartures(selectedStop Stop) []Departure {
	resp, err := http.DefaultClient.Get(baseUrl + "/departureBoard?format=json&id=" + selectedStop.ID)
	if err != nil {
		panic(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var response DepartureResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		panic(err)
	}

	return response.DepartureBoard.Departure
}
