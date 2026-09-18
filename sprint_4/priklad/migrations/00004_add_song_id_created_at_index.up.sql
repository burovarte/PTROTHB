CREATE INDEX idx_playlist_songs_song_id_created_at
ON public.playlist_songs(song_id, created_at);

DROP INDEX public.idx_playlist_songs_song_id;