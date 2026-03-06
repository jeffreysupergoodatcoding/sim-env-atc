package adsb

import (
	"time"
)

// DemoAircraft defines the hardcoded aircraft for the replay demo
var DemoAircraft = []*Aircraft{
	{
		Hex:      "SIM120",
		Flight:   "SIM120",
		Airline:  "Cessna",
		Status:   "grounded",
		OnGround: true,
		LastSeen: time.Now(),
		ADSB: &ADSBTarget{
			Hex:         "SIM120",
			Flight:      "SIM120",
			Type:        "sim",
			Lat:         37.4611,
			Lon:         -122.1150,
			AltBaro:     0,
			TrueHeading: 298,
			TAS:         0,
		},
	},
	{
		Hex:      "SIM201",
		Flight:   "SIM201",
		Airline:  "Piper",
		Status:   "grounded",
		OnGround: true,
		LastSeen: time.Now(),
		ADSB: &ADSBTarget{
			Hex:         "SIM201",
			Flight:      "SIM201",
			Type:        "sim",
			Lat:         37.4650,
			Lon:         -122.1100,
			AltBaro:     0,
			TrueHeading: 130,
			TAS:         0,
		},
	},
	{
		Hex:      "SIM305",
		Flight:   "SIM305",
		Airline:  "Cessna",
		Status:   "airborne",
		OnGround: false,
		LastSeen: time.Now(),
		ADSB: &ADSBTarget{
			Hex:         "SIM305",
			Flight:      "SIM305",
			Type:        "sim",
			Lat:         37.4580,
			Lon:         -122.1200,
			AltBaro:     1500,
			TrueHeading: 270,
			TAS:         120,
		},
	},
}
