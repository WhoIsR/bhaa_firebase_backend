package models

import "gorm.io/gorm"

// Cetakan buat tabel transaksi utama
type Order struct {
	gorm.Model
	UserID     uint    `json:"user_id"`
	TotalPrice float64 `json:"total_price"`
	Status     string  `json:"status"`
	// Sambungan ke rincian barang (Satu pesanan bisa banyak barang)
	OrderItems []OrderItem `gorm:"foreignKey:OrderID" json:"items"`
}

// Cetakan buat tabel rincian barang yang dibeli
type OrderItem struct {
	gorm.Model
	OrderID   uint    `json:"order_id"`
	ProductID uint    `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}
