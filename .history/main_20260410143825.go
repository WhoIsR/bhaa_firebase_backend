package main

import (
	"log"
	"os"

	"github.com/WhoIsR/bhaa_firebase_backend/config"
	"github.com/WhoIsR/bhaa_firebase_backend/routes"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Buka .env
	if err := godotenv.Load(); err != nil {
		log.Println("File .env tidak ditemukan, menggunakan environment variable sistem")
	}

	// 2. Nyalain koneksi ke Google Firebase
	config.InitFirebase()

	// 3. Nyalain koneksi ke MySQL Laragon dan otomatis bikin tabel
	config.InitDatabase()

	// 4. Gelar peta jalan (Routes)
	router := routes.SetupRouter()

	// 5. Hidupkan Server di port 8080
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server berjalan di http://localhost:%s", port)
	log.Printf("Health check: http://localhost:%s/v1/health", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
