package http

import (
	"encoding/json"
	"errors"
	nethttp "net/http"
	"time"

	"github/burovarte/PTROTHB/sprint_4/priklad/internal/domain"
)

type PlaybackService interface {
	AddSong(song domain.Song) error
	Play() error
	Pause() error
	Next() error
	Prev() error
}

type PlaybackHandler struct {
	service PlaybackService
}

func NewPlaybackHandler(service PlaybackService) *PlaybackHandler {
	return &PlaybackHandler{service: service}
}

type addPlaybackSongRequest struct {
	Name       string `json:"name"`
	DurationMS int64  `json:"duration_ms"`
}

func (h *PlaybackHandler) AddSong(w nethttp.ResponseWriter, r *nethttp.Request) {
	var request addPlaybackSongRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		nethttp.Error(w, "invalid request body", nethttp.StatusBadRequest)
		return
	}

	song := domain.Song{
		Name:     request.Name,
		Duration: time.Duration(request.DurationMS) * time.Millisecond,
	}

	if err := h.service.AddSong(song); errors.Is(err, domain.ErrInvalidSong) {
		nethttp.Error(w, "invalid song", nethttp.StatusBadRequest)
		return
	} else if err != nil {
		nethttp.Error(w, "internal server error", nethttp.StatusInternalServerError)
		return
	}

	w.WriteHeader(nethttp.StatusNoContent)
}

func (h *PlaybackHandler) Play(w nethttp.ResponseWriter, _ *nethttp.Request) {
	err := h.service.Play()
	if errors.Is(err, domain.ErrEmptyPlaylist) {
		nethttp.Error(w, "playlist is empty", nethttp.StatusConflict)
		return
	}
	if err != nil {
		nethttp.Error(w, "internal server error", nethttp.StatusInternalServerError)
		return
	}

	w.WriteHeader(nethttp.StatusNoContent)
}

func (h *PlaybackHandler) Pause(w nethttp.ResponseWriter, _ *nethttp.Request) {
	err := h.service.Pause()
	if errors.Is(err, domain.ErrNotPlaying) {
		nethttp.Error(w, "playlist is not playing", nethttp.StatusConflict)
		return
	}
	if err != nil {
		nethttp.Error(w, "internal server error", nethttp.StatusInternalServerError)
		return
	}

	w.WriteHeader(nethttp.StatusNoContent)
}

func (h *PlaybackHandler) Next(w nethttp.ResponseWriter, _ *nethttp.Request) {
	err := h.service.Next()
	if errors.Is(err, domain.ErrNoCurrentSong) || errors.Is(err, domain.ErrNoNextSong) {
		nethttp.Error(w, "next song is not available", nethttp.StatusConflict)
		return
	}
	if err != nil {
		nethttp.Error(w, "internal server error", nethttp.StatusInternalServerError)
		return
	}

	w.WriteHeader(nethttp.StatusNoContent)
}

func (h *PlaybackHandler) Prev(w nethttp.ResponseWriter, _ *nethttp.Request) {
	err := h.service.Prev()
	if errors.Is(err, domain.ErrNoCurrentSong) || errors.Is(err, domain.ErrNoPrevSong) {
		nethttp.Error(w, "previous song is not available", nethttp.StatusConflict)
		return
	}
	if err != nil {
		nethttp.Error(w, "internal server error", nethttp.StatusInternalServerError)
		return
	}

	w.WriteHeader(nethttp.StatusNoContent)
}
