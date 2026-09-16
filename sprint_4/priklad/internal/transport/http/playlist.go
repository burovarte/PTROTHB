package http

import (
	"context"
	"encoding/json"
	"errors"
	nethttp "net/http"
	"strconv"
	"time"

	"github/burovarte/PTROTHB/sprint_4/priklad/internal/domain"
)

type PlaylistService interface {
	CreatePlaylist(ctx context.Context, playlist domain.Playlist) (domain.Playlist, error)
	ReadPlaylist(ctx context.Context, id int64) (domain.Playlist, error)
	UpdatePlaylist(ctx context.Context, playlist domain.Playlist) error
	DeletePlaylist(ctx context.Context, id int64) error
	AddSongToPlaylist(ctx context.Context, playlistID int64, songID int64, position int64) error
}

type PlaylistHandler struct {
	service PlaylistService
}

func NewPlaylistHandler(service PlaylistService) *PlaylistHandler {
	return &PlaylistHandler{service: service}
}

type createPlaylistRequest struct {
	Name string `json:"name"`
}

type updatePlaylistRequest struct {
	Name string `json:"name"`
}

type addSongToPlaylistRequest struct {
	SongID   int64 `json:"song_id"`
	Position int64 `json:"position"`
}

type playlistResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (h *PlaylistHandler) CreatePlaylist(w nethttp.ResponseWriter, r *nethttp.Request) {
	var request createPlaylistRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		nethttp.Error(w, "invalid request body", nethttp.StatusBadRequest)
		return
	}

	createdPlaylist, err := h.service.CreatePlaylist(r.Context(), domain.Playlist{Name: request.Name})
	if errors.Is(err, domain.ErrInvalidPlaylist) {
		nethttp.Error(w, "invalid playlist", nethttp.StatusBadRequest)
		return
	}
	if err != nil {
		nethttp.Error(w, "internal server error", nethttp.StatusInternalServerError)
		return
	}

	response := newPlaylistResponse(createdPlaylist)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(nethttp.StatusCreated)
	_ = json.NewEncoder(w).Encode(response)
}

func (h *PlaylistHandler) ReadPlaylist(w nethttp.ResponseWriter, r *nethttp.Request) {
	id, ok := playlistIDFromRequest(w, r)
	if !ok {
		return
	}

	playlist, err := h.service.ReadPlaylist(r.Context(), id)
	if errors.Is(err, domain.ErrPlaylistNotFound) {
		nethttp.Error(w, "playlist not found", nethttp.StatusNotFound)
		return
	}
	if err != nil {
		nethttp.Error(w, "internal server error", nethttp.StatusInternalServerError)
		return
	}

	response := newPlaylistResponse(playlist)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(nethttp.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func (h *PlaylistHandler) UpdatePlaylist(w nethttp.ResponseWriter, r *nethttp.Request) {
	id, ok := playlistIDFromRequest(w, r)
	if !ok {
		return
	}

	var request updatePlaylistRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		nethttp.Error(w, "invalid request body", nethttp.StatusBadRequest)
		return
	}

	err := h.service.UpdatePlaylist(r.Context(), domain.Playlist{ID: id, Name: request.Name})
	if errors.Is(err, domain.ErrInvalidPlaylist) {
		nethttp.Error(w, "invalid playlist", nethttp.StatusBadRequest)
		return
	}
	if errors.Is(err, domain.ErrPlaylistNotFound) {
		nethttp.Error(w, "playlist not found", nethttp.StatusNotFound)
		return
	}
	if err != nil {
		nethttp.Error(w, "internal server error", nethttp.StatusInternalServerError)
		return
	}

	w.WriteHeader(nethttp.StatusNoContent)
}

func (h *PlaylistHandler) DeletePlaylist(w nethttp.ResponseWriter, r *nethttp.Request) {
	id, ok := playlistIDFromRequest(w, r)
	if !ok {
		return
	}

	err := h.service.DeletePlaylist(r.Context(), id)
	if errors.Is(err, domain.ErrPlaylistNotFound) {
		nethttp.Error(w, "playlist not found", nethttp.StatusNotFound)
		return
	}
	if err != nil {
		nethttp.Error(w, "internal server error", nethttp.StatusInternalServerError)
		return
	}

	w.WriteHeader(nethttp.StatusNoContent)
}

func (h *PlaylistHandler) AddSongToPlaylist(w nethttp.ResponseWriter, r *nethttp.Request) {
	playlistID, ok := playlistIDFromRequest(w, r)
	if !ok {
		return
	}

	var request addSongToPlaylistRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		nethttp.Error(w, "invalid request body", nethttp.StatusBadRequest)
		return
	}

	err := h.service.AddSongToPlaylist(r.Context(), playlistID, request.SongID, request.Position)
	if errors.Is(err, domain.ErrInvalidID) || errors.Is(err, domain.ErrInvalidPosition) {
		nethttp.Error(w, "invalid song id or position", nethttp.StatusBadRequest)
		return
	}
	if err != nil {
		nethttp.Error(w, "internal server error", nethttp.StatusInternalServerError)
		return
	}

	w.WriteHeader(nethttp.StatusNoContent)
}

func playlistIDFromRequest(w nethttp.ResponseWriter, r *nethttp.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		nethttp.Error(w, "invalid playlist id", nethttp.StatusBadRequest)
		return 0, false
	}

	return id, true
}

func newPlaylistResponse(playlist domain.Playlist) playlistResponse {
	return playlistResponse{
		ID:        playlist.ID,
		Name:      playlist.Name,
		CreatedAt: playlist.CreatedAt,
		UpdatedAt: playlist.UpdatedAt,
	}
}
