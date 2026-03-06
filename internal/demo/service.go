package demo

import (
	"sync"

	"github.com/yegors/co-atc/internal/simulation"
	"github.com/yegors/co-atc/pkg/logger"
)

// Service manages the replay demo sequence
type Service struct {
	clips      []Clip
	currentIdx int
	simService *simulation.Service
	logger     *logger.Logger
	mu         sync.Mutex
}

// NewService creates a new demo service
func NewService(simService *simulation.Service, logger *logger.Logger) *Service {
	return &Service{
		clips:      GetDefaultClips(),
		currentIdx: -1,
		simService: simService,
		logger:     logger.Named("demo"),
	}
}

// GetClips returns all clips in the sequence
func (s *Service) GetClips() []Clip {
	return s.clips
}

// GetCurrentClip returns the current clip info
func (s *Service) GetCurrentClip() (Clip, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.currentIdx < 0 || s.currentIdx >= len(s.clips) {
		return Clip{}, false
	}
	return s.clips[s.currentIdx], true
}

// NextClip advances to the next clip and returns it
func (s *Service) NextClip() (Clip, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.currentIdx++
	if s.currentIdx >= len(s.clips) {
		s.currentIdx = 0 // loop for now or stop
	}

	clip := s.clips[s.currentIdx]

	// Trigger simulation action if necessary
	if clip.Callsign != "" {
		s.logger.Info("Triggering simulation action for demo",
			logger.String("callsign", clip.Callsign),
			logger.String("action", clip.Action))

		// Ensure aircraft exists in simulation with the right callsign
		// We use callsign as hex for simplicity in the demo
		hex := clip.Callsign
		_, err := s.simService.RegisterAircraft(hex, clip.Callsign, clip.TargetLat, clip.TargetLon, clip.TargetAlt, clip.TargetHdg, clip.TargetSpd, 0)
		if err != nil {
			s.logger.Warn("Failed to register aircraft for demo", logger.Error(err))
		}

		err = s.simService.UpdateControls(hex, clip.TargetHdg, clip.TargetSpd, 0)
		if err != nil {
			s.logger.Warn("Failed to update simulation controls for demo", logger.Error(err))
		}
	}

	return clip, true
}

// Reset resets the demo sequence
func (s *Service) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.currentIdx = -1
}
