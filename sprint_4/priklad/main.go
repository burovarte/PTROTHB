package priklad

import (
	"errors"
	"time"
)

type playbackState int

var ErrEmptyPlaylist = errors.New("playlist is empty")

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
	head    *node
	current *node
	tail    *node
	state   playbackState
}

type Player interface {
	Play() error
	Pause() error
	AddSong(song Song) error
	Next() error
	Prev() error
}

func (p *Playlist) Play() error {
	if p.head == nil {
		return ErrEmptyPlaylist
	}

	return nil
}
