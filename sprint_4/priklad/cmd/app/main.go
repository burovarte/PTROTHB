package main

import (
	"context"
	"database/sql"
	"log"
	nethttp "net/http"
	"os"
	"time"

	"github/burovarte/PTROTHB/sprint_4/priklad/internal/domain"
	"github/burovarte/PTROTHB/sprint_4/priklad/internal/repository/postgres"
	"github/burovarte/PTROTHB/sprint_4/priklad/internal/service"
	httptransport "github/burovarte/PTROTHB/sprint_4/priklad/internal/transport/http"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	db, err := sql.Open("pgx", databaseURL)

	if err != nil {
		log.Fatalf("open database: %v", err)
	}

	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	err = db.PingContext(ctx)

	if err != nil {
		log.Fatalf("ping database: %v", err)
	}

	log.Println("database connection established")

	songRepository := postgres.NewSongRepository(db)
	songService := service.NewSongService(songRepository)
	songHandler := httptransport.NewSongHandler(songService)

	playlistRepository := postgres.NewPlaylistRepository(db)
	playlistService := service.NewPlaylistService(playlistRepository)
	playlistHandler := httptransport.NewPlaylistHandler(playlistService)

	player := &domain.Playback{}
	playbackService := service.NewPlaybackService(player)
	playbackHandler := httptransport.NewPlaybackHandler(playbackService)

	router := httptransport.NewRouter(songHandler, playlistHandler, playbackHandler)

	server := &nethttp.Server{
		Addr:              ":8080",
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("HTTP server started on :8080")

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("HTTP server: %v", err)
	}
}
