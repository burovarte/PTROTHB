package priklad

import (
	"errors"
	"sync"
	"testing"
	"time"
)

func TestPlaylistPlay_EmptyPlaylist(t *testing.T) {
	playlist := Playlist{}
	err := playlist.Play()
	if !errors.Is(err, ErrEmptyPlaylist) {
		t.Fatalf("expected ErrEmptyPlaylist, got %v", err)
	}
	if playlist.current != nil {
		t.Error("expected current to be nil")
	}
	if playlist.state != stateStopped {
		t.Errorf("expected stateStopped, got %v", playlist.state)
	}
}

func TestPlaylistAddSong(t *testing.T) {
	playlist := Playlist{}
	song1 := Song{Name: "песня 1", Duration: time.Minute}
	song2 := Song{Name: "песня 2", Duration: 2 * time.Minute}
	if err := playlist.AddSong(song1); err != nil {
		t.Fatalf("add first song: %v", err)
	}
	if err := playlist.AddSong(song2); err != nil {
		t.Fatalf("add second song: %v", err)
	}
	if playlist.head == nil || playlist.tail == nil {
		t.Fatal("expected non-nil head and tail")
	}
	if playlist.head.song != song1 {
		t.Errorf("expected first song %+v, got %+v", song1, playlist.head.song)
	}
	if playlist.tail.song != song2 {
		t.Errorf("expected second song %+v, got %+v", song2, playlist.tail.song)
	}
	if playlist.head.next != playlist.tail {
		t.Error("expected head.next to point to tail")
	}
	if playlist.tail.prev != playlist.head {
		t.Error("expected tail.prev to point to head")
	}
	if playlist.current != nil {
		t.Error("expected current to remain nil")
	}
}

func TestPlaylistAddSong_InvalidDuration(t *testing.T) {
	playlist := Playlist{}
	err := playlist.AddSong(Song{Name: "invalid"})
	if !errors.Is(err, ErrInvalidSong) {
		t.Fatalf("expected ErrInvalidSong, got %v", err)
	}
	if playlist.head != nil || playlist.tail != nil {
		t.Error("invalid song must not change the list")
	}
}

func TestPlaylistPlay_StartsHead(t *testing.T) {
	playlist := Playlist{}
	if err := playlist.AddSong(Song{Name: "песня", Duration: time.Minute}); err != nil {
		t.Fatalf("add song: %v", err)
	}
	err := playlist.Play()
	t.Cleanup(func() { _ = playlist.Pause() })
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if playlist.current != playlist.head {
		t.Error("expected current to point to head")
	}
	if playlist.state != statePlaying {
		t.Errorf("expected statePlaying, got %v", playlist.state)
	}
	if playlist.cancel == nil {
		t.Error("expected active cancel function")
	}
}

func TestPlaylistPlay_WhilePlayingIsNoOp(t *testing.T) {
	playlist := Playlist{}
	if err := playlist.AddSong(Song{Name: "песня", Duration: time.Minute}); err != nil {
		t.Fatalf("add song: %v", err)
	}
	if err := playlist.Play(); err != nil {
		t.Fatalf("first play: %v", err)
	}
	t.Cleanup(func() { _ = playlist.Pause() })
	startedAt := playlist.startedAt
	if err := playlist.Play(); err != nil {
		t.Fatalf("second play: %v", err)
	}
	if !playlist.startedAt.Equal(startedAt) {
		t.Error("repeated Play must not restart playback")
	}
}

func TestPlaylistPause_NotPlaying(t *testing.T) {
	playlist := Playlist{}
	if err := playlist.Pause(); !errors.Is(err, ErrNotPlaying) {
		t.Fatalf("expected ErrNotPlaying, got %v", err)
	}
}

func TestPlaylistPauseAndResume(t *testing.T) {
	playlist := Playlist{}
	duration := 500 * time.Millisecond
	if err := playlist.AddSong(Song{Name: "песня", Duration: duration}); err != nil {
		t.Fatalf("add song: %v", err)
	}
	if err := playlist.Play(); err != nil {
		t.Fatalf("play: %v", err)
	}
	time.Sleep(30 * time.Millisecond)
	if err := playlist.Pause(); err != nil {
		t.Fatalf("pause: %v", err)
	}
	pausedRemaining := playlist.remaining
	if playlist.state != statePaused {
		t.Errorf("expected statePaused, got %v", playlist.state)
	}
	if pausedRemaining <= 0 || pausedRemaining >= duration {
		t.Errorf("expected remaining between 0 and %v, got %v", duration, pausedRemaining)
	}
	if err := playlist.Play(); err != nil {
		t.Fatalf("resume: %v", err)
	}
	t.Cleanup(func() { _ = playlist.Pause() })
	if playlist.state != statePlaying {
		t.Errorf("expected statePlaying after resume, got %v", playlist.state)
	}
	if playlist.remaining != pausedRemaining {
		t.Errorf("expected resume from %v, got %v", pausedRemaining, playlist.remaining)
	}
}

func TestPlaylistNext(t *testing.T) {
	playlist := Playlist{}
	first := Song{Name: "первая", Duration: 100 * time.Millisecond}
	second := Song{Name: "вторая", Duration: time.Second}
	if err := playlist.AddSong(first); err != nil {
		t.Fatalf("add first song: %v", err)
	}
	if err := playlist.AddSong(second); err != nil {
		t.Fatalf("add second song: %v", err)
	}
	if err := playlist.Play(); err != nil {
		t.Fatalf("play: %v", err)
	}
	if err := playlist.Next(); err != nil {
		t.Fatalf("next: %v", err)
	}
	t.Cleanup(func() { _ = playlist.Pause() })
	if playlist.current != playlist.tail {
		t.Error("expected current to point to the second song")
	}
	if playlist.state != statePlaying {
		t.Errorf("expected statePlaying, got %v", playlist.state)
	}
	if playlist.remaining != second.Duration {
		t.Errorf("expected duration %v, got %v", second.Duration, playlist.remaining)
	}
	time.Sleep(150 * time.Millisecond)
	playlist.mu.Lock()
	current, state := playlist.current, playlist.state
	playlist.mu.Unlock()
	if current != playlist.tail || state != statePlaying {
		t.Error("old playback must not change the new current song")
	}
}

func TestPlaylistNext_Errors(t *testing.T) {
	t.Run("no current song", func(t *testing.T) {
		playlist := Playlist{}
		if err := playlist.Next(); !errors.Is(err, ErrNoCurrentSong) {
			t.Fatalf("expected ErrNoCurrentSong, got %v", err)
		}
	})
	t.Run("no next song", func(t *testing.T) {
		playlist := Playlist{}
		if err := playlist.AddSong(Song{Name: "единственная", Duration: time.Minute}); err != nil {
			t.Fatalf("add song: %v", err)
		}
		if err := playlist.Play(); err != nil {
			t.Fatalf("play: %v", err)
		}
		t.Cleanup(func() { _ = playlist.Pause() })
		if err := playlist.Next(); !errors.Is(err, ErrNoNextSong) {
			t.Fatalf("expected ErrNoNextSong, got %v", err)
		}
	})
}

func TestPlaylistPrev(t *testing.T) {
	playlist := Playlist{}
	for _, song := range []Song{
		{Name: "первая", Duration: time.Second},
		{Name: "вторая", Duration: time.Second},
	} {
		if err := playlist.AddSong(song); err != nil {
			t.Fatalf("add song: %v", err)
		}
	}
	if err := playlist.Play(); err != nil {
		t.Fatalf("play: %v", err)
	}
	if err := playlist.Next(); err != nil {
		t.Fatalf("next: %v", err)
	}
	if err := playlist.Prev(); err != nil {
		t.Fatalf("prev: %v", err)
	}
	t.Cleanup(func() { _ = playlist.Pause() })
	if playlist.current != playlist.head {
		t.Error("expected current to return to the first song")
	}
	if err := playlist.Prev(); !errors.Is(err, ErrNoPrevSong) {
		t.Fatalf("expected ErrNoPrevSong, got %v", err)
	}
}

func TestPlaylist_AutoAdvanceAndStop(t *testing.T) {
	playlist := Playlist{}
	for _, name := range []string{"первая", "вторая"} {
		if err := playlist.AddSong(Song{Name: name, Duration: 20 * time.Millisecond}); err != nil {
			t.Fatalf("add song: %v", err)
		}
	}
	if err := playlist.Play(); err != nil {
		t.Fatalf("play: %v", err)
	}
	waitFor(t, time.Second, func() bool {
		playlist.mu.Lock()
		defer playlist.mu.Unlock()
		return playlist.current == playlist.tail && playlist.state == statePlaying
	})
	waitFor(t, time.Second, func() bool {
		playlist.mu.Lock()
		defer playlist.mu.Unlock()
		return playlist.current == playlist.tail && playlist.state == stateStopped
	})
}

func TestPlaylistAddSong_Concurrent(t *testing.T) {
	playlist := Playlist{}
	const songCount = 100
	var wg sync.WaitGroup
	errs := make(chan error, songCount)
	for i := 0; i < songCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errs <- playlist.AddSong(Song{Name: "песня", Duration: time.Second})
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatalf("concurrent AddSong: %v", err)
		}
	}
	playlist.mu.Lock()
	defer playlist.mu.Unlock()
	count := 0
	var previous *node
	for current := playlist.head; current != nil; current = current.next {
		if current.prev != previous {
			t.Fatal("broken prev link after concurrent AddSong")
		}
		previous = current
		count++
	}
	if count != songCount {
		t.Fatalf("expected %d songs, got %d", songCount, count)
	}
	if previous != playlist.tail {
		t.Error("tail must point to the final node")
	}
}

func TestPlaylistNextPrev_Concurrent(t *testing.T) {
	playlist := Playlist{}
	for _, name := range []string{"первая", "вторая", "третья"} {
		if err := playlist.AddSong(Song{Name: name, Duration: time.Minute}); err != nil {
			t.Fatalf("add song: %v", err)
		}
	}
	if err := playlist.Play(); err != nil {
		t.Fatalf("play: %v", err)
	}
	if err := playlist.Next(); err != nil {
		t.Fatalf("move to the middle song: %v", err)
	}
	t.Cleanup(func() { _ = playlist.Pause() })

	const operationCount = 100
	var wg sync.WaitGroup
	unexpectedErrors := make(chan error, operationCount)

	for i := 0; i < operationCount; i++ {
		wg.Add(1)
		go func(moveNext bool) {
			defer wg.Done()
			if moveNext {
				if err := playlist.Next(); err != nil && !errors.Is(err, ErrNoNextSong) {
					unexpectedErrors <- err
				}
				return
			}
			if err := playlist.Prev(); err != nil && !errors.Is(err, ErrNoPrevSong) {
				unexpectedErrors <- err
			}
		}(i%2 == 0)
	}

	wg.Wait()
	close(unexpectedErrors)
	for err := range unexpectedErrors {
		t.Errorf("unexpected Next/Prev error: %v", err)
	}

	playlist.mu.Lock()
	defer playlist.mu.Unlock()
	if playlist.current == nil {
		t.Fatal("current song must not become nil")
	}
	if playlist.head.prev != nil {
		t.Error("head.prev must remain nil")
	}
	if playlist.tail.next != nil {
		t.Error("tail.next must remain nil")
	}
	if playlist.head.next == nil || playlist.head.next.prev != playlist.head {
		t.Error("forward and backward links must remain consistent")
	}
	if playlist.tail.prev == nil || playlist.tail.prev.next != playlist.tail {
		t.Error("backward and forward links must remain consistent")
	}
}

func waitFor(t *testing.T, timeout time.Duration, condition func() bool) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if condition() {
			return
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatal("condition was not met before timeout")
}
