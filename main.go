package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"url-shortener/config"
	"url-shortener/database"
	"url-shortener/repository"
	"url-shortener/routes"
	"url-shortener/services"
	"url-shortener/handlers"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	db, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	urlRepo := repository.NewURLRepository(db)
	urlService := services.NewURLService(urlRepo)
	urlHandler := handlers.NewURLHandler(urlService,cfg.BaseURL)


	defer db.Close()

	log.Println("Connected to PostgreSQL successfully!")

	mux := http.NewServeMux()
	routes.Register(mux, urlHandler)
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Println("Server running on port 8080")

	log.Fatal(server.ListenAndServe())
}
