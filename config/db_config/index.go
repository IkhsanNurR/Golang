package db_config

import (
	"log"
	"os"
)

var DB_DRIVER = ""
var DB_HOST = ""
var DB_PORT = ""
var DB_USER = ""
var DB_PASSWORD = ""
var DB_NAME = ""
var DB_SCHEMA = ""

func InitDbConfig() {
	if os.Getenv("DB_DRIVER") != "postgres" {
		log.Fatal("DB_DRIVER is not set")
	} else {
		DB_DRIVER = os.Getenv("DB_DRIVER")
		DB_HOST = os.Getenv("DB_HOST")
		DB_PORT = os.Getenv("DB_PORT")
		DB_USER = os.Getenv("DB_USER")
		DB_PASSWORD = os.Getenv("DB_PASSWORD")
		DB_NAME = os.Getenv("DB_NAME")
		DB_SCHEMA = os.Getenv("DB_SCHEMA")
	}
}
