package main

import (
	"log"

	"github.com/WhoIsR/bhaa_firebase_backend/config"
	"github.com/WhoIsR/bhaa_firebase_backend/models"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	config.InitDatabase()

	products := []models.Product{
		{Name: "Nasi Goreng Spesial", Price: 25000, Category: "Makanan", Stock: 50, Description: "Nasi goreng dengan telur dan ayam", ImageURL: "https://picsum.photos/400"},
		{Name: "Sate Ayam 10 Tusuk", Price: 20000, Category: "Makanan", Stock: 100, Description: "Sate ayam dengan bumbu kacang", ImageURL: "https://picsum.photos/401"},
		{Name: "Nasi Padang", Price: 22000, Category: "Makanan", Stock: 60, Description: "Nasi padang dengan lauk dan sambal", ImageURL: "https://picsum.photos/402"},
		{Name: "Nasi Kuning", Price: 18000, Category: "Makanan", Stock: 40, Description: "Nasi kuning dengan lauk dan sambal", ImageURL: "https://picsum.photos/403"},
		{Name: "Es Teh Manis", Price: 8000, Category: "Minuman", Stock: 200, Description: "Es teh manis segar", ImageURL: "https://picsum.photos/404"},
		{Name: "Kopi Susu", Price: 15000, Category: "Minuman", Stock: 150, Description: "Kopi susu kekinian", ImageURL: "https://picsum.photos/405"},
		{Name: "Ayam Bakar", Price: 30000, Category: "Makanan", Stock: 30, Description: "Ayam bakar dengan sambal", ImageURL: "https://picsum.photos/406"},
		{Name: "Jus Alpukat", Price: 12000, Category: "Minuman", Stock: 80, Description: "Jus alpukat segar", ImageURL: "https://picsum.photos/407"},
		{Name: "Nasi Uduk", Price: 20000, Category: "Makanan", Stock: 70, Description: "Nasi uduk dengan lauk dan sambal", ImageURL: "https://picsum.photos/408"},
		{Name: "Es Campur", Price: 10000, Category: "Minuman", Stock: 120, Description: "Es campur segar dengan berbagai topping", ImageURL: "https://picsum.photos/409"},
	}

	for _, p := range products {
		config.DB.Create(&p)
	}
	log.Printf("Seed berhasil: %d produk ditambahkan", len(products))
}
