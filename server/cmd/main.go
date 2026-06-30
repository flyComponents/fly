package main

import (
	"fly_server/internal/discovery"
	internal_http "fly_server/internal/http"
	"fly_server/internal/logger"
	"fly_server/internal/service"
	"fly_server/internal/storage"
	"net/http"
)

// @title			Fly Server API
// @version		1.0
// @description	Drone swarm control API
// @host
// @BasePath		/89.104.67.132
func main() {
	log := logger.New()

	repo := storage.NewMemory()             // метаданные
	fs := storage.NewFSStorage("./uploads") // файлы

	photoService := service.NewPhotoService(repo, fs)
	handler := internal_http.NewHandler(repo, log)
	photoHandler := internal_http.NewPhotoHandler(photoService)

	discovery.Start(log, 37020, 8080)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: internal_http.NewRouter(handler, photoHandler),
	}

	log.Println("server started")
	log.Fatal(srv.ListenAndServe())
}
