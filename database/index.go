package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"main.go/config/db_config"
)

var DB *gorm.DB

func ConnectDatabase() (*gorm.DB, error) {
	var err error
	log.Default().Println("konek database coba")
	log.Default().Println(db_config.DB_DRIVER)
	if db_config.DB_DRIVER == "postgres" {
		log.Default().Println("masuk")
		dsn := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			db_config.DB_HOST, db_config.DB_PORT, db_config.DB_USER, db_config.DB_PASSWORD, db_config.DB_NAME,
		)
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				TablePrefix: db_config.DB_SCHEMA + ".",
			},
		})

		if err != nil {
			log.Fatalf("Error connecting to database: %v", err)
		} else {
			log.Default().Println("Database connected")
		}
	}
	return DB, nil
}

// migrate -database "postgres://postgres:postgres@localhost:5432/Asep?sslmode=disable" -path database/migrations up
