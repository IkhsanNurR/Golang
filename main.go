package main

import (
	"log"

	"github.com/joho/godotenv"
	"main.go/bootstrap"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	bootstrap.BootstrapApp()
}
