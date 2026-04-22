package repositories

import (
	"fmt"

	"github.com/WhoIsR/bhaa_firebase_backend/models"
	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrder(order *models.Order) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db}
}

// Fungsi nyimpen ke database MySQL + kurangi stok produk
func (r *orderRepository) CreateOrder(order *models.Order) error {
	// Pakai Transaction biar aman: kalau stok gagal dikurangi, pesanan juga dibatalkan
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Simpan pesanan utama + item-itemnya
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		// 2. Kurangi stok setiap produk yang dibeli
		for _, item := range order.OrderItems {
			result := tx.Model(&models.Product{}).
				Where("id = ? AND stock >= ?", item.ProductID, item.Quantity).
				Update("stock", gorm.Expr("stock - ?", item.Quantity))

			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return fmt.Errorf("stok produk ID %d tidak mencukupi", item.ProductID)
			}
		}

		return nil
	})
}
