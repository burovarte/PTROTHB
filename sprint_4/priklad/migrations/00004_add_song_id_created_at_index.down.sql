DROP INDEX idx_playlist_songs_song_id_created_at;

CREATE INDEX idx_playlist_songs_song_id
ON public.playlist_songs (song_id);