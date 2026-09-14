//go:build integration

package postgres

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"testing"
	"time"

	"github/burovarte/PTROTHB/sprint_4/priklad/internal/domain"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestSongRepositoryCreateAndRead(t *testing.T) {
	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL == "" {
		t.Fatal("DATABASE_URL environment variable is not set")
	}

	db, err := sql.Open("pgx", databaseURL)

	if err != nil {
		t.Fatalf("open database: %v", err)
	}

	t.Cleanup(func() {
		_ = db.Close()
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	err = db.PingContext(ctx)

	if err != nil {
		t.Fatalf("ping database: %v", err)
	}

	repository := NewSongRepository(db)

	inputSong := domain.Song{
		Name:     "lalal",
		Duration: 1500 * time.Millisecond,
	}

	createdSong, err := repository.CreateSong(ctx, inputSong)
	if err != nil {
		t.Fatalf("create song: %v", err)
	}

	if createdSong.ID <= 0 {
		t.Errorf("created song ID = %d, want a positive value", createdSong.ID)
	}
	if createdSong.Name != inputSong.Name {
		t.Errorf("created song name = %q, want %q", createdSong.Name, inputSong.Name)
	}
	if createdSong.Duration != inputSong.Duration {
		t.Errorf("created song duration = %v, want %v", createdSong.Duration, inputSong.Duration)
	}
	if createdSong.CreatedAt.IsZero() {
		t.Error("created song CreatedAt is zero")
	}
	if createdSong.UpdatedAt.IsZero() {
		t.Error("created song UpdatedAt is zero")
	}

	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cleanupCancel()

		if _, cleanupErr := db.ExecContext(cleanupCtx, "DELETE FROM public.song WHERE id = $1", createdSong.ID); cleanupErr != nil {
			t.Errorf("cleanup song %d: %v", createdSong.ID, cleanupErr)
		}
	})

	readSong, err := repository.ReadSong(ctx, createdSong.ID)
	if err != nil {
		t.Fatalf("read song: %v", err)
	}

	if readSong.ID != createdSong.ID {
		t.Errorf("read song ID = %d, want %d", readSong.ID, createdSong.ID)
	}
	if readSong.Name != createdSong.Name {
		t.Errorf("read song name = %q, want %q", readSong.Name, createdSong.Name)
	}
	if readSong.Duration != createdSong.Duration {
		t.Errorf("read song duration = %v, want %v", readSong.Duration, createdSong.Duration)
	}
	if !readSong.CreatedAt.Equal(createdSong.CreatedAt) {
		t.Errorf("read song CreatedAt = %v, want %v", readSong.CreatedAt, createdSong.CreatedAt)
	}
	if !readSong.UpdatedAt.Equal(createdSong.UpdatedAt) {
		t.Errorf("read song UpdatedAt = %v, want %v", readSong.UpdatedAt, createdSong.UpdatedAt)
	}

	canceledCtx, cancelImmediately := context.WithCancel(context.Background())
	cancelImmediately()

	if _, err = repository.ReadSong(canceledCtx, createdSong.ID); !errors.Is(err, context.Canceled) {
		t.Errorf("ReadSong() with canceled context error = %v, want context.Canceled", err)
	}

	deleteResult, err := db.ExecContext(ctx, "DELETE FROM public.song WHERE id = $1", createdSong.ID)
	if err != nil {
		t.Fatalf("delete created song: %v", err)
	}

	deletedRows, err := deleteResult.RowsAffected()
	if err != nil {
		t.Fatalf("get deleted row count: %v", err)
	}
	if deletedRows != 1 {
		t.Fatalf("deleted row count = %d, want 1", deletedRows)
	}

	if _, err = repository.ReadSong(ctx, createdSong.ID); !errors.Is(err, domain.ErrSongNotFound) {
		t.Errorf("ReadSong() after delete error = %v, want %v", err, domain.ErrSongNotFound)
	}

}
