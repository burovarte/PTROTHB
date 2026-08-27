package priklad

import (
	"errors"
	"testing"
)

func TestPlaylistPlay_EmptyPlaylist(t *testing.T) {
	playlist := Playlist{}

	err := playlist.Play()

	if !errors.Is(err, ErrEmptyPlaylist) {
		t.Fatalf("expected ErrEmptyPlaylist, got %v", err)
	}

	if playlist.current != nil {
		t.Errorf("expected current to be nil")
	}

	if playlist.state != stateStopped {
		t.Errorf("expected state is not stopped ")
	}
}
