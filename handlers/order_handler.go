package handlers

import (
	"net/http"

	"github.com/WhoIsR/bhaa_firebase_backend/models"
	"github.com/WhoIsR/bhaa_firebase_backend/services"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	service services.OrderService
}

func NewOrderHandler(service services.OrderService) *OrderHandler {
	return &OrderHandler{service}
}

func (h *OrderHandler) Checkout(c *gin.Context) {
	// 1. Bikin wadah kosong buat nangkep JSON kiriman Flutter
	var req struct {
		TotalPrice float64 `json:"total_price" binding:"required"`
		Items      []struct {
			ProductID uint    `json:"product_id" binding:"required"`
			Quantity  int     `json:"quantity" binding:"required"`
			Price     float64 `json:"price" binding:"required"`
		} `json:"items" binding:"required"`
	}

	// 2. Masukin data JSON ke dalam wadah tadi
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data belanjaan nggak lengkap bro"})
		return
	}

	// 3. Ambil ID user dari Satpam Middleware (yang didapet dari token JWT)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Lu belum login bro"})
		return
	}

	// Konversi tipe data userID biar aman
	var finalUserID uint
	switch v := userID.(type) {
	case float64:
		finalUserID = uint(v)
	case uint:
		finalUserID = v
	case int:
		finalUserID = uint(v)
	}

	// 4. Susun data jadi bentuk Model Database
	order := models.Order{
		UserID:     finalUserID,
		TotalPrice: req.TotalPrice,
		Status:     "success", // Simulasi sukses sesuai permintaan dosen
	}

	for _, item := range req.Items {
		order.OrderItems = append(order.OrderItems, models.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		})
	}

	// 5. Suruh Mandor (Service) nyimpen datanya
	if err := h.service.CreateOrder(&order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal nyimpen pesanan ke database"})
		return
	}

	// 6. Kasih struk sukses ke Flutter
	c.JSON(http.StatusCreated, gin.H{"message": "Checkout berhasil masuk MySQL!", "order": order})
}
