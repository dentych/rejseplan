package main

import (
	"github.com/dentych/rejseplan/rejseplan"
	"github.com/pterm/pterm"
)

func main() {
	input, err := pterm.DefaultInteractiveTextInput.Show("Find stop")
	if err != nil {
		panic(err)
	}

	stops := rejseplan.GetStops(input)
	var stopNames []string
	for _, stop := range stops {
		stopNames = append(stopNames, stop.Name)
	}

	var selectedStop rejseplan.StopLocation
	if len(stops) == 0 {
		pterm.DefaultBasicText.Println("No stops found")
		return
	} else if len(stops) == 1 {
		selectedStop = stops[0]
	} else {
		name, err := pterm.DefaultInteractiveSelect.WithOptions(stopNames).Show("Select a stop")
		if err != nil {
			panic(err)
		}

		for _, stop := range stops {
			if stop.Name == name {
				selectedStop = stop
				break
			}
		}
	}

	pterm.DefaultBasicText.Printf("You have chosen stop '%s' with ID '%s'\n", selectedStop.Name, selectedStop.ID)

	departures := rejseplan.GetDepartures(selectedStop)

	data := pterm.TableData{
		{"Line", "Direction", "Time"},
	}

	for _, departure := range departures {
		data = append(data, []string{departure.Line, departure.Direction, getTime(departure.ActualTime, departure.PlannedTime)})
	}

	pterm.DefaultTable.WithData(data).Render()
}

func getTime(text, def string) string {
	if text != "" {
		return text
	}

	return def
}
