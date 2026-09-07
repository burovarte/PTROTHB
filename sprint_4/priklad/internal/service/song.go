package service

import (
	"context"

	"github/burovarte/PTROTHB/sprint_4/priklad/internal/domain"
)

type SongRepository interface {
	CreateSong(ctx context.Context, song domain.Song) (domain.Song, error)
	ReadSong(ctx context.Context, id int64) (domain.Song, error)
	UpdateSong(ctx context.Context, song domain.Song) error
	DeleteSong(ctx context.Context, id int64) error
}

type SongService struct {
	repo SongRepository
}

func NewSongService(repo SongRepository) *SongService {
	service := SongService{
		repo: repo,
	}

	return &service
}

func (s *SongService) CreateSong(ctx context.Context, song domain.Song) (domain.Song, error) {
	if song.Duration <= 0 || song.Name == "" {
		return domain.Song{}, domain.ErrInvalidSong
	}

	return s.repo.CreateSong(ctx, song)
}

func (s *SongService) ReadSong(ctx context.Context, id int64) (domain.Song, error) {
	if id <= 0 {
		return domain.Song{}, domain.ErrInvalidID
	}

	return s.repo.ReadSong(ctx, id)
}

func (s *SongService) UpdateSong(ctx context.Context, song domain.Song) error {
	if song.ID <= 0 {
		return domain.ErrInvalidID
	}

	if song.Duration <= 0 || song.Name == "" {
		return domain.ErrInvalidSong
	}

	return s.repo.UpdateSong(ctx, song)
}

func (s *SongService) DeleteSong(ctx context.Context, id int64) error {
	if id <= 0 {
		return domain.ErrInvalidID
	}

	return s.repo.DeleteSong(ctx, id)
}
