package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Hajdudev/ecoDatabase/internal/app"
	"github.com/Hajdudev/ecoDatabase/internal/routes"
)

func main() {
	log.Print("starting server...")

	application, err := app.NewApplication()
	if err != nil {
		log.Fatalf("failed to start application: %v", err)
	}

	application.Logger.Println("We are running the app")

	r := routes.SetupRoutes(application)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3003"
		log.Printf("defaulting to port %s", port)
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Printf("listening on port %s", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
