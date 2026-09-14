package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github/burovarte/PTROTHB/sprint_4/priklad/internal/domain"
)

type SongRepository struct {
	db *sql.DB
}

func NewSongRepository(db *sql.DB) *SongRepository {
	return &SongRepository{db: db}
}
func (r *SongRepository) CreateSong(ctx context.Context, song domain.Song) (domain.Song, error) {
	const query = `INSERT INTO public.song(name, duration)
	VALUES($1, $2)
	RETURNING id, name, duration, created_at, updated_at`

	row := r.db.QueryRowContext(ctx, query, song.Name, int64(song.Duration))

	var createdSong domain.Song
	var duration int64

	err := row.Scan(&createdSong.ID,
		&createdSong.Name,
		&duration,
		&createdSong.CreatedAt,
		&createdSong.UpdatedAt)

	if err != nil {
		return domain.Song{}, fmt.Errorf("create song: %w", err)
	}

	createdSong.Duration = time.Duration(duration)

	return createdSong, nil
}

func (r *SongRepository) ReadSong(ctx context.Context, id int64) (domain.Song, error) {
	const query = `
		SELECT id, name, duration, created_at, updated_at
		FROM public.song
		WHERE id = $1
	`

	row := r.db.QueryRowContext(ctx, query, id)

	var song domain.Song
	var duration int64

	err := row.Scan(&song.ID, &song.Name, &duration, &song.CreatedAt, &song.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return domain.Song{}, domain.ErrSongNotFound
	}

	if err != nil {
		return domain.Song{}, fmt.Errorf("read song: %w", err)
	}

	song.Duration = time.Duration(duration)

	return song, nil
}
