package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github/burovarte/PTROTHB/sprint_4/priklad/internal/domain"
)

type PlaylistRepository struct {
	db *sql.DB
}

func NewPlaylistRepository(db *sql.DB) *PlaylistRepository {
	return &PlaylistRepository{db: db}
}

func (r *PlaylistRepository) ReadPlaylist(ctx context.Context, id int64) (domain.Playlist, error) {
	const query = `
		SELECT id, name, created_at, updated_at
		FROM public.playlists
		WHERE id = $1
	`

	row := r.db.QueryRowContext(ctx, query, id)

	var playlist domain.Playlist

	err := row.Scan(&playlist.ID, &playlist.Name, &playlist.CreatedAt, &playlist.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return domain.Playlist{}, domain.ErrPlaylistNotFound
	}

	if err != nil {
		return domain.Playlist{}, fmt.Errorf("read playlist: %w", err)
	}

	return playlist, nil
}

func (r *PlaylistRepository) CreatePlaylist(ctx context.Context, playlist domain.Playlist) (domain.Playlist, error) {
	const query = `INSERT INTO public.playlists(name)
	VALUES($1)
	RETURNING id, name, created_at, updated_at`

	row := r.db.QueryRowContext(ctx, query, playlist.Name)

	var createdPlaylist domain.Playlist

	err := row.Scan(&createdPlaylist.ID,
		&createdPlaylist.Name,
		&createdPlaylist.CreatedAt,
		&createdPlaylist.UpdatedAt)

	if err != nil {
		return domain.Playlist{}, fmt.Errorf("create playlist: %w", err)
	}

	return createdPlaylist, nil
}

func (r *PlaylistRepository) UpdatePlaylist(ctx context.Context, playlist domain.Playlist) error {
	const query = `UPDATE public.playlists
	SET name = $1, updated_at = NOW()
	WHERE id = $2
	`

	res, err := r.db.ExecContext(ctx, query, playlist.Name, playlist.ID)

	if err != nil {
		return fmt.Errorf("update playlist: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("update playlist: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrPlaylistNotFound
	}

	return nil
}

func (r *PlaylistRepository) DeletePlaylist(ctx context.Context, id int64) error {
	const query = `DELETE FROM public.playlists
	WHERE id = $1
	`

	res, err := r.db.ExecContext(ctx, query, id)

	if err != nil {
		return fmt.Errorf("delete playlist: %w", err)
	}

	resAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("get deleted row count: %w", err)
	}

	if resAffected == 0 {
		return domain.ErrPlaylistNotFound
	}

	return nil
}

func (r *PlaylistRepository) AddSongToPlaylist(ctx context.Context, playlistID int64, songID int64, position int64) error {
	const query = `
		INSERT INTO public.playlist_songs(playlist_id, song_id, position)
		VALUES ($1, $2, $3)
	`
	const countQuery = `
		SELECT COUNT(*)
		FROM public.playlist_songs
		WHERE playlist_id = $1
	`

	const shiftQuery = `
		UPDATE public.playlist_songs
		SET position = position + 1,
			updated_at = NOW()
		WHERE playlist_id = $1
			AND position >= $2
	`

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin add song transaction: %w", err)
	}

	defer tx.Rollback()

	var songsCount int64
	if err = tx.QueryRowContext(ctx, countQuery, playlistID).Scan(&songsCount); err != nil {
		return fmt.Errorf("count playlist songs: %w", err)
	}

	if position > songsCount+1 {
		return domain.ErrInvalidPosition
	}

	_, err = tx.ExecContext(ctx, shiftQuery, playlistID, position)
	if err != nil {
		return fmt.Errorf("shift playlist songs: %w", err)
	}

	_, err = tx.ExecContext(ctx, query, playlistID, songID, position)
	if err != nil {
		return fmt.Errorf("add song to playlist: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit add song transaction: %w", err)
	}

	return nil

}
