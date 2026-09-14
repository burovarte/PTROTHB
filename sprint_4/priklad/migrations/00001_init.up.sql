CREATE TABLE "public".song (
    id bigint NOT NULL GENERATED ALWAYS AS IDENTITY,
    name text NOT NULL,
    duration bigint NOT NULL,
    updated_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT pk_song PRIMARY KEY(id)
);


ALTER TABLE "public".song ADD CONSTRAINT chk_song_duration CHECK ( duration > 0 );


CREATE TABLE "public".playlists(
    id bigint NOT NULL GENERATED ALWAYS AS IDENTITY,
    name text NOT NULL,
    updated_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT pk_playlists PRIMARY KEY(id)
);

CREATE TABLE "public".playlist_songs(
    id bigint NOT NULL GENERATED ALWAYS AS IDENTITY,
    CONSTRAINT pk_playlist_songs PRIMARY KEY(id),
    playlist_id bigint NOT NULL,
    CONSTRAINT fk_playlist_songs_playlist FOREIGN KEY (playlist_id) REFERENCES public.playlists(id),
    song_id bigint NOT NULL,
    CONSTRAINT fk_playlist_songs_id FOREIGN KEY (song_id) REFERENCES public.song(id),
    position bigint NOT NULL CHECK (position > 0),
    CONSTRAINT unq_playlist_songs_position UNIQUE (playlist_id, position ),
    updated_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL
);