package app_config

import (
	"log"
	"os"
)

var APP_PORT = ":8080"

func InitAppConfig() {

	if os.Getenv("APP_PORT") == "" {
		log.Default().Println("APP_PORT is not set. Defaulting to :8080")
	} else {
		APP_PORT = os.Getenv("APP_PORT")
	}
}
