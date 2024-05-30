package main

import (
	"log"

	"github.com/joho/godotenv"
	"gorm.io/gen"
	"main.go/config"
	"main.go/database"
)

func main() {
	err := godotenv.Load("../.env")
	log.Default().Println(err)
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	config.InitConfig()

	// Memanggil fungsi generate
	generate()
}

func generate() {
	// Menghubungkan ke database
	log.Println("Koneksi ke database...")
	DB, err := database.ConnectDatabase()

	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	// Membuat generator dengan konfigurasi tertentu
	g := gen.NewGenerator(gen.Config{
		OutPath: "../generated/models",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})
	g.UseDB(DB)

	log.Println("Generate model dari database...")

	// Mengaplikasikan fungsi-fungsi dasar untuk menghasilkan tabel-tabel
	g.ApplyBasic(g.GenerateAllTable()...)

	// Menjalankan generator
	g.Execute()
}
