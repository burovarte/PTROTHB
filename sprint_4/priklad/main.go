package priklad

import (
	"context"
	"errors"
	"sync"
	"time"
)

type playbackState int

var ErrEmptyPlaylist = errors.New("playlist is empty")

var ErrNotPlaying = errors.New("playlist is not playing")

var ErrNoCurrentSong = errors.New("has not current")

var ErrNoNextSong = errors.New("no more songs")

var ErrNoPrevSong = errors.New("no prev songs")

var ErrInvalidSong = errors.New("invalid song")

const (
	stateStopped playbackState = iota
	statePlaying
	statePaused
)

type Song struct {
	Name     string
	Duration time.Duration
}

type node struct {
	prev *node
	song Song
	next *node
}

type Playlist struct {
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

func (p *Playlist) Play() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.head == nil {
		return ErrEmptyPlaylist
	}

	if p.current == nil {
		p.current = p.head
	}

	p.state = statePlaying

	return nil
}

func (p *Playlist) AddSong(song Song) error {
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

func (p *Playlist) Pause() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.state != statePlaying {
		return ErrNotPlaying

	} else {
		p.state = statePaused
	}

	return nil
}

func (p *Playlist) Next() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.current == nil {
		return ErrNoCurrentSong
	}

	if p.current.next == nil {
		return ErrNoNextSong
	}

	p.current = p.current.next

	p.state = statePlaying

	return nil
}

func (p *Playlist) Prev() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.current == nil {
		return ErrNoCurrentSong
	}

	if p.current.prev == nil {
		return ErrNoPrevSong
	}

	p.current = p.current.prev

	p.state = statePlaying

	return nil
}
