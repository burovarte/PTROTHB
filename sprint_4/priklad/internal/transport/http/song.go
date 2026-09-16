package http

import (
	"context"
	"encoding/json"
	"errors"
	"github/burovarte/PTROTHB/sprint_4/priklad/internal/domain"
	nethttp "net/http"
	"strconv"
	"time"
)

type SongService interface {
	CreateSong(ctx context.Context, song domain.Song) (domain.Song, error)
	ReadSong(ctx context.Context, id int64) (domain.Song, error)
	UpdateSong(ctx context.Context, song domain.Song) error
	DeleteSong(ctx context.Context, id int64) error
}

type SongHandler struct {
	service SongService
}

func NewSongHandler(service SongService) *SongHandler {
	return &SongHandler{
		service: service,
	}
}

type createSongRequest struct {
	Name       string `json:"name"`
	DurationMS int64  `json:"duration_ms"`
}

type updateSongRequest struct {
	Name       string `json:"name"`
	DurationMS int64  `json:"duration_ms"`
}

type songResponse struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	DurationMS int64     `json:"duration_ms"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (h *SongHandler) CreateSong(w nethttp.ResponseWriter, r *nethttp.Request) {
	var request createSongRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		nethttp.Error(w, "invalid request body", nethttp.StatusBadRequest)
		return
	}

	song := domain.Song{
		Name:     request.Name,
		Duration: time.Duration(request.DurationMS) * time.Millisecond,
	}

	createdSong, err := h.service.CreateSong(r.Context(), song)

	if errors.Is(err, domain.ErrInvalidSong) {
		nethttp.Error(w, "invalid song", nethttp.StatusBadRequest)
		return
	}

	if err != nil {
		nethttp.Error(w, "internal server error", nethttp.StatusInternalServerError)
		return
	}

	response := songResponse{
		ID:         createdSong.ID,
		Name:       createdSong.Name,
		DurationMS: createdSong.Duration.Milliseconds(),
		CreatedAt:  createdSong.CreatedAt,
		UpdatedAt:  createdSong.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(nethttp.StatusCreated)
	_ = json.NewEncoder(w).Encode(response)
}

func (h *SongHandler) ReadSong(w nethttp.ResponseWriter, r *nethttp.Request) {
	idText := r.PathValue("id")

	id, err := strconv.ParseInt(idText, 10, 64)
	if err != nil || id <= 0 {
		nethttp.Error(w, "invalid song id", nethttp.StatusBadRequest)
		return
	}

	song, err := h.service.ReadSong(r.Context(), id)
	if errors.Is(err, domain.ErrSongNotFound) {
		nethttp.Error(w, "song not found", nethttp.StatusNotFound)
		return
	}

	if err != nil {
		nethttp.Error(w, "internal server error", nethttp.StatusInternalServerError)
		return
	}

	response := songResponse{
		ID:         song.ID,
		Name:       song.Name,
		DurationMS: song.Duration.Milliseconds(),
		CreatedAt:  song.CreatedAt,
		UpdatedAt:  song.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(nethttp.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func (h *SongHandler) UpdateSong(w nethttp.ResponseWriter, r *nethttp.Request) {
	idText := r.PathValue("id")

	id, err := strconv.ParseInt(idText, 10, 64)

	if err != nil || id <= 0 {
		nethttp.Error(w, "invalid song id", nethttp.StatusBadRequest)
		return
	}

	var request updateSongRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		nethttp.Error(w, "invalid request body", nethttp.StatusBadRequest)
		return
	}

	song := domain.Song{
		ID:       id,
		Name:     request.Name,
		Duration: time.Duration(request.DurationMS) * time.Millisecond,
	}

	err = h.service.UpdateSong(r.Context(), song)

	if errors.Is(err, domain.ErrInvalidSong) {
		nethttp.Error(w, "invalid song", nethttp.StatusBadRequest)
		return
	}

	if errors.Is(err, domain.ErrSongNotFound) {
		nethttp.Error(w, "song not found", nethttp.StatusNotFound)
		return
	}

	if err != nil {
		nethttp.Error(w, "internal server error", nethttp.StatusInternalServerError)
		return
	}

	w.WriteHeader(nethttp.StatusNoContent)

}

func (h *SongHandler) DeleteSong(w nethttp.ResponseWriter, r *nethttp.Request) {
	idText := r.PathValue("id")

	id, err := strconv.ParseInt(idText, 10, 64)
	if err != nil || id <= 0 {
		nethttp.Error(w, "invalid song id", nethttp.StatusBadRequest)
		return
	}

	err = h.service.DeleteSong(r.Context(), id)

	if errors.Is(err, domain.ErrSongNotFound) {
		nethttp.Error(w, "song not found", nethttp.StatusNotFound)
		return
	}

	if err != nil {
		nethttp.Error(w, "internal server error", nethttp.StatusInternalServerError)
		return
	}

	w.WriteHeader(nethttp.StatusNoContent)

}
