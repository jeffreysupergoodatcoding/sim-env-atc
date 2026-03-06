package models

import "time"

// TranscriptionRecord represents a transcription record in the database
type TranscriptionRecord struct {
	ID               int64     `json:"id"`
	FrequencyID      string    `json:"frequency_id"`
	CreatedAt        time.Time `json:"created_at"`
	Content          string    `json:"content"`
	IsComplete       bool      `json:"is_complete"`
	IsProcessed      bool      `json:"is_processed"`
	ContentProcessed string    `json:"content_processed"`
	SpeakerType      string    `json:"speaker_type,omitempty"` // "ATC" or "PILOT"
	Callsign         string    `json:"callsign,omitempty"`     // Aircraft callsign if speaker is a pilot
}

// ChatMessage represents a message in a chat session
type ChatMessage struct {
	ID        string    `json:"id"`
	SessionID string    `json:"session_id"`
	Type      string    `json:"type"` // "user", "assistant", "system"
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
	AudioData []byte    `json:"audio_data,omitempty"`
}
