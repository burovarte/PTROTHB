package service

import (
	"context"

	"github/burovarte/PTROTHB/sprint_4/priklad/internal/domain"
)

type PlaylistRepository interface {
	CreatePlaylist(ctx context.Context, playList domain.Playlist) (domain.Playlist, error)
	ReadPlaylist(ctx context.Context, id int64) (domain.Playlist, error)
	UpdatePlaylist(ctx context.Context, playList domain.Playlist) error
	DeletePlaylist(ctx context.Context, id int64) error
	AddSongToPlaylist(ctx context.Context, songID int64, playlistID int64, position int64) error
}

type PlaylistService struct {
	repo PlaylistRepository
}

func NewPlaylistService(repo PlaylistRepository) *PlaylistService {
	service := PlaylistService{
		repo: repo,
	}

	return &service
}

func (p *PlaylistService) CreatePlaylist(ctx context.Context, playlist domain.Playlist) (domain.Playlist, error) {
	if playlist.Name == "" {
		return domain.Playlist{}, domain.ErrInvalidPlaylist
	}

	return p.repo.CreatePlaylist(ctx, playlist)
}

func (p *PlaylistService) ReadPlaylist(ctx context.Context, id int64) (domain.Playlist, error) {
	if id <= 0 {
		return domain.Playlist{}, domain.ErrInvalidID
	}

	return p.repo.ReadPlaylist(ctx, id)
}

func (p *PlaylistService) UpdatePlaylist(ctx context.Context, playlist domain.Playlist) error {
	if playlist.ID <= 0 {
		return domain.ErrInvalidID
	}

	if playlist.Name == "" {
		return domain.ErrInvalidPlaylist
	}

	return p.repo.UpdatePlaylist(ctx, playlist)
}

func (p *PlaylistService) DeletePlaylist(ctx context.Context, id int64) error {
	if id <= 0 {
		return domain.ErrInvalidID
	}

	return p.repo.DeletePlaylist(ctx, id)
}
