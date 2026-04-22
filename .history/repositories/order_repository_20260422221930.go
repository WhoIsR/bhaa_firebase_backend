package repositories

import (
	"bhaa_firebase_backend/models"

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

// Fungsi nyimpen ke database MySQL
func (r *orderRepository) CreateOrder(order *models.Order) error {
	return r.db.Create(order).Error
}
