package main

import (
	"github.com/joho/godotenv"
	"github.com/lukasbischof/luk4s.dev/app"
	"log"
	"os"
)

func main() {
	if os.Getenv("APP_ENV") == "development" {
		log.Println("Starting in development mode")
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	app.Run()
}
