package main

import (
	"context"
	"log"
	"net/http"
	"time"
	"url-shortener/config"
	"url-shortener/database"
	"url-shortener/handlers"
	"url-shortener/redis"
	"url-shortener/repository"
	"url-shortener/routes"
	"url-shortener/services"
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

	rdb , err := redis.Connect(ctx)
	if err != nil {
		log.Fatal("Redis connection failed:" , err)
	}

	defer rdb.Close()


	urlRepo := repository.NewURLRepository(db)
	urlCache := repository.NewCacheRepository(rdb)
	urlService := services.NewURLService(urlRepo,urlCache)
	urlHandler := handlers.NewURLHandler(urlService,cfg.BaseURL)


	defer db.Close()

	log.Println("Connected to PostgreSQL successfully!")
	log.Println("Connected to Redis successfully!")

	mux := http.NewServeMux()
	routes.Register(mux, urlHandler)
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Println("Server running on port 8080")

	log.Fatal(server.ListenAndServe())
}
