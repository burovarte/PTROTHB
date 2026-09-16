package http

import (
	nethttp "net/http"
)

func NewRouter(
	songHandler *SongHandler,
	playlistHandler *PlaylistHandler,
	playbackHandler *PlaybackHandler,
) nethttp.Handler {

	mux := nethttp.NewServeMux()

	mux.HandleFunc("POST /songs", songHandler.CreateSong)
	mux.HandleFunc("GET /songs/{id}", songHandler.ReadSong)
	mux.HandleFunc("PUT /songs/{id}", songHandler.UpdateSong)
	mux.HandleFunc("DELETE /songs/{id}", songHandler.DeleteSong)

	mux.HandleFunc("POST /playlists", playlistHandler.CreatePlaylist)
	mux.HandleFunc("GET /playlists/{id}", playlistHandler.ReadPlaylist)
	mux.HandleFunc("PUT /playlists/{id}", playlistHandler.UpdatePlaylist)
	mux.HandleFunc("DELETE /playlists/{id}", playlistHandler.DeletePlaylist)
	mux.HandleFunc("POST /playlists/{id}/songs", playlistHandler.AddSongToPlaylist)

	mux.HandleFunc("POST /playback/songs", playbackHandler.AddSong)
	mux.HandleFunc("POST /playback/play", playbackHandler.Play)
	mux.HandleFunc("POST /playback/pause", playbackHandler.Pause)
	mux.HandleFunc("POST /playback/next", playbackHandler.Next)
	mux.HandleFunc("POST /playback/prev", playbackHandler.Prev)

	return mux
}
