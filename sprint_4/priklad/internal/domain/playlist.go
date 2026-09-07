package domain

import (
	"context"
	"sync"
	"time"
)

type playbackState int

const (
	stateStopped playbackState = iota
	statePlaying
	statePaused
)

type node struct {
	prev *node
	song Song
	next *node
}

type Playlist struct {
	ID        int64
	Name      string
	UpdatedAt time.Time
	CreatedAt time.Time
}

type Playback struct {
	head      *node
	current   *node
	tail      *node
	state     playbackState
	mu        sync.Mutex
	startedAt time.Time
	remaining time.Duration
	cancel    context.CancelFunc
}

type Player interface {
	Play() error
	Pause() error
	AddSong(song Song) error
	Next() error
	Prev() error
}

func (p *Playback) Play() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.head == nil {
		return ErrEmptyPlaylist
	}

	if p.current == nil {
		p.current = p.head
	}

	if p.state == statePlaying {
		return nil
	}

	duration := p.current.song.Duration

	if p.state == statePaused && p.remaining > 0 {
		duration = p.remaining
	}

	p.startPlaybackLocked(duration)

	return nil
}

func (p *Playback) AddSong(song Song) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if song.Duration <= 0 {
		return ErrInvalidSong
	}
	newNode := &node{
		song: song,
	}

	if p.head == nil {
		p.head = newNode
		p.tail = newNode
		return nil
	} else {
		newNode.prev = p.tail
		p.tail.next = newNode
		p.tail = newNode

		return nil
	}

}

func (p *Playback) Pause() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.state != statePlaying {
		return ErrNotPlaying

	}

	if p.cancel != nil {
		p.cancel()
	}

	p.cancel = nil

	timePassed := time.Since(p.startedAt)

	p.remaining -= timePassed

	if p.remaining < 0 {
		p.remaining = 0
	}

	p.state = statePaused

	return nil
}

func (p *Playback) Next() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.current == nil {
		return ErrNoCurrentSong
	}

	if p.current.next == nil {
		return ErrNoNextSong
	}

	if p.cancel != nil {
		p.cancel()
	}

	p.cancel = nil

	p.current = p.current.next

	p.startPlaybackLocked(p.current.song.Duration)

	return nil
}

func (p *Playback) Prev() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.current == nil {
		return ErrNoCurrentSong
	}

	if p.current.prev == nil {
		return ErrNoPrevSong
	}

	if p.cancel != nil {
		p.cancel()
	}

	p.cancel = nil

	p.current = p.current.prev

	p.startPlaybackLocked(p.current.song.Duration)

	return nil
}

func (p *Playback) runPlayback(ctx context.Context, duration time.Duration, track *node) {
	timer := time.NewTimer(duration)

	defer timer.Stop()

	select {
	case <-ctx.Done():

	case <-timer.C:
		p.mu.Lock()
		defer p.mu.Unlock()

		if ctx.Err() != nil || p.current != track {
			return
		}

		p.cancel = nil
		p.remaining = 0

		if p.current.next == nil {
			p.state = stateStopped
			return
		}

		p.current = track.next
		p.startPlaybackLocked(p.current.song.Duration)

	}
}

func (p *Playback) startPlaybackLocked(duration time.Duration) {

	ctx, cancel := context.WithCancel(context.Background())

	p.cancel = cancel

	p.startedAt = time.Now()

	p.remaining = duration

	p.state = statePlaying

	go p.runPlayback(ctx, duration, p.current)
}
