package services

import (
	"github.com/WhoIsR/bhaa_firebase_backend/models"
	"github.com/WhoIsR/bhaa_firebase_backend/repositories"
)

type OrderService interface {
	CreateOrder(order *models.Order) error
	GetOrdersByUserID(userID uint, page, limit int) ([]models.Order, error)
}

type orderService struct {
	repo repositories.OrderRepository
}

func NewOrderService(repo repositories.OrderRepository) OrderService {
	return &orderService{repo}
}

// Fungsi logika sebelum dikirim ke database
func (s *orderService) CreateOrder(order *models.Order) error {
	return s.repo.CreateOrder(order)
}

func (s *orderService) GetOrdersByUserID(userID uint, page, limit int) ([]models.Order, error) {
	return s.repo.GetOrdersByUserID(userID, page, limit)
}
