package service

import "github/burovarte/PTROTHB/sprint_4/priklad/internal/domain"

type PlaybackService struct {
	player domain.Player
}

func NewPlaybackService(player domain.Player) *PlaybackService {
	return &PlaybackService{player: player}
}

func (s *PlaybackService) Play() error {
	return s.player.Play()
}

func (s *PlaybackService) Pause() error {
	return s.player.Pause()
}

func (s *PlaybackService) Next() error {
	return s.player.Next()
}

func (s *PlaybackService) Prev() error {
	return s.player.Prev()
}

func (s *PlaybackService) AddSong(song domain.Song) error {
	return s.player.AddSong(song)
}
