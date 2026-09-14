package domain

import "errors"

var (
	ErrEmptyPlaylist = errors.New("playlist is empty")

	ErrNotPlaying = errors.New("playlist is not playing")

	ErrNoCurrentSong = errors.New("has not current")

	ErrNoNextSong = errors.New("no more songs")

	ErrNoPrevSong = errors.New("no prev songs")

	ErrInvalidSong = errors.New("invalid song")

	ErrInvalidID = errors.New("invalid id")

	ErrInvalidPlaylist = errors.New("invalid playlist")

	ErrInvalidPosition = errors.New("invalid position")

	ErrSongNotFound = errors.New("song not found")
)
