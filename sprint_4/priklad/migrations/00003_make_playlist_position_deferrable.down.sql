ALTER TABLE public.playlist_songs
DROP CONSTRAINT unq_playlist_songs_position;

ALTER TABLE public.playlist_songs
ADD CONSTRAINT unq_playlist_songs_position
UNIQUE (playlist_id, position);