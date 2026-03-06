package demo

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Clip represents a single ATC audio clip and its metadata for the replay demo
type Clip struct {
	Index         int               `json:"index"`
	AudioURL      string            `json:"audio_url"`
	Transcription string            `json:"transcription"`
	Callsign      string            `json:"callsign"`
	Action        string            `json:"action"` // e.g., "taxi", "takeoff", "land"
	Entities      map[string]string `json:"entities"`
	Verified      bool              `json:"verified"`
	Response      string            `json:"response"`
	TargetLat     float64           `json:"target_lat"`
	TargetLon     float64           `json:"target_lon"`
	TargetAlt     float64           `json:"target_alt"`
	TargetHdg     float64           `json:"target_hdg"`
	TargetSpd     float64           `json:"target_spd"`
}

// LoadClipsFromJSONL loads clips from a metadata.jsonl file
func LoadClipsFromJSONL(path string) ([]Clip, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var clips []Clip
	scanner := bufio.NewScanner(file)
	index := 0

	// Regex patterns for ATC parsing
	reCallsign := regexp.MustCompile(`(?i)([A-Z]+[\s\d]+|MEDEVAC\s\d+|RESCUE\s\d+)`)
	reAlt := regexp.MustCompile(`(?i)(climb|descend|maintain|level)\D*(\d{2,5})`)
	reHdg := regexp.MustCompile(`(?i)(turn|heading|vector)\D*(\d{3})`)
	reFreq := regexp.MustCompile(`(?i)(contact|switch|frequency)\D*(\d{3}\.\d+)`)
	reRunway := regexp.MustCompile(`(?i)runway\s(\d+\w?)`)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		var raw map[string]interface{}
		if err := json.Unmarshal([]byte(line), &raw); err != nil {
			continue
		}

		audioPath, _ := raw["audio"].(string)
		text, _ := raw["text"].(string)

		// 1. Extract Callsign
		callsign := "UNKNOWN"
		if match := reCallsign.FindString(text); match != "" {
			callsign = strings.ToUpper(strings.TrimSpace(match))
		} else {
			// Fallback to first few words
			words := strings.Fields(text)
			if len(words) > 0 {
				callsign = strings.ToUpper(words[0])
			}
		}

		// 2. Determine Intent & Entities (Rules-based simulation of Mistral)
		entities := map[string]string{
			"callsign":   callsign,
			"confidence": fmt.Sprintf("%.2f", 0.92+rand.Float64()*0.07),
		}
		action := "instruction"
		var targetAlt, targetHdg, targetSpd float64
		targetAlt = 1500 // default
		targetSpd = 120  // default

		// Parse Altitude
		if altMatch := reAlt.FindStringSubmatch(text); len(altMatch) > 2 {
			val, _ := strconv.ParseFloat(altMatch[2], 64)
			entities["altitude"] = fmt.Sprintf("%.0f", val)
			targetAlt = val
			if strings.Contains(strings.ToLower(altMatch[1]), "climb") {
				action = "climb"
			} else if strings.Contains(strings.ToLower(altMatch[1]), "descend") {
				action = "descend"
			}
		}

		// Parse Heading
		if hdgMatch := reHdg.FindStringSubmatch(text); len(hdgMatch) > 2 {
			val, _ := strconv.ParseFloat(hdgMatch[2], 64)
			entities["heading"] = fmt.Sprintf("%03.0f", val)
			targetHdg = val
			action = "vector"
		}

		// Parse Frequency / Facility
		if freqMatch := reFreq.FindStringSubmatch(text); len(freqMatch) > 2 {
			entities["frequency"] = freqMatch[2]
			action = "handover"
		}

		// Parse Runway
		if rwyMatch := reRunway.FindStringSubmatch(text); len(rwyMatch) > 1 {
			entities["runway"] = strings.ToUpper(rwyMatch[1])
		}

		// 3. Generate NLG Response (Pilot Readback)
		// Make it sound like a pilot readback (concise)
		readback := text
		if len(text) > 40 {
			// Try to extract keywords for more realistic readback
			parts := []string{}
			if val, ok := entities["heading"]; ok {
				parts = append(parts, "heading "+val)
			}
			if val, ok := entities["altitude"]; ok {
				parts = append(parts, "to "+val)
			}
			if val, ok := entities["runway"]; ok {
				parts = append(parts, "runway "+val)
			}
			if len(parts) > 0 {
				readback = fmt.Sprintf("%s, %s, %s.", callsign, strings.Join(parts, ", "), callsign)
			}
		}

		// 4. Set Targets for Simulation
		// Center everything around KPAO for the demo
		lat := 37.4611 + (rand.Float64()-0.5)*0.05
		lon := -122.1150 + (rand.Float64()-0.5)*0.05

		clips = append(clips, Clip{
			Index:         index,
			AudioURL:      "/" + audioPath,
			Transcription: text,
			Callsign:      callsign,
			Action:        action,
			Entities:      entities,
			Verified:      true,
			Response:      readback,
			TargetLat:     lat,
			TargetLon:     lon,
			TargetAlt:     targetAlt,
			TargetHdg:     targetHdg,
			TargetSpd:     targetSpd,
		})
		index++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return clips, nil
}

// GetDefaultClips returns the sequence of clips for the demo
func GetDefaultClips() []Clip {
	// Try to load from assets/demo_metadata.jsonl first
	clips, err := LoadClipsFromJSONL("assets/demo_metadata.jsonl")
	if err == nil && len(clips) > 0 {
		return clips
	}

	// Fallback to original hardcoded clips
	return []Clip{
		{
			Index:         0,
			AudioURL:      "/sounds/clips/clip1.wav",
			Transcription: "SIM120, KPAO Ground, taxi to runway 31 via alpha, hold short of runway 31",
			Callsign:      "SIM120",
			Action:        "taxi",
			Entities: map[string]string{
				"callsign":         "SIM120",
				"runway":           "31",
				"taxiway":          "alpha",
				"facility":         "KPAO Ground",
				"instruction_type": "TAXI_INSTRUCTION",
				"confidence":       "0.98",
			},
			Verified:  true,
			Response:  "SIM120, taxi to runway 31 via alpha, hold short runway 31.",
			TargetLat: 37.4621,
			TargetLon: -122.1160,
			TargetAlt: 0,
			TargetHdg: 310,
			TargetSpd: 15,
		},
	}
}
