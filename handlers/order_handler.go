package handlers

import (
	"net/http"

	"github.com/WhoIsR/bhaa_firebase_backend/config"
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
		ShippingAddress string  `json:"shipping_address"`
		Notes           string  `json:"notes"`
		PaymentMethod   string  `json:"payment_method"`
		TotalPrice      float64 `json:"total_price"`
		Items           []struct {
			ProductID uint    `json:"product_id"`
			Quantity  int     `json:"quantity"`
			Price     float64 `json:"price"`
		} `json:"items"`
	}

	// 2. Masukin data JSON ke dalam wadah tadi
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Data checkout tidak valid"})
		return
	}

	// 3. Ambil ID user dari Satpam Middleware (yang didapet dari token JWT)
	finalUserID, ok := userIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "User belum login"})
		return
	}

	if len(req.Items) == 0 {
		var cartItems []models.CartItem
		if err := config.DB.Preload("Product").Where("user_id = ?", finalUserID).Find(&cartItems).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengambil cart"})
			return
		}
		if len(cartItems) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Cart masih kosong"})
			return
		}
		for _, item := range cartItems {
			req.Items = append(req.Items, struct {
				ProductID uint    `json:"product_id"`
				Quantity  int     `json:"quantity"`
				Price     float64 `json:"price"`
			}{
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Price:     item.Product.Price,
			})
			req.TotalPrice += item.Product.Price * float64(item.Quantity)
		}
	}

	// 4. Susun data jadi bentuk Model Database
	order := models.Order{
		UserID:     finalUserID,
		TotalPrice: req.TotalPrice,
		Status:     "success",
	}

	for _, item := range req.Items {
		if item.ProductID == 0 || item.Quantity <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Data item checkout tidak valid"})
			return
		}
		order.OrderItems = append(order.OrderItems, models.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		})
	}

	// 5. Suruh Mandor (Service) nyimpen datanya
	if err := h.service.CreateOrder(&order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal menyimpan pesanan"})
		return
	}

	config.DB.Where("user_id = ?", finalUserID).Delete(&models.CartItem{})

	// 6. Kasih struk sukses ke Flutter
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Checkout berhasil",
		"data": gin.H{
			"id":               order.ID,
			"total_amount":     order.TotalPrice,
			"status":           order.Status,
			"shipping_address": req.ShippingAddress,
			"notes":            req.Notes,
			"payment_method":   req.PaymentMethod,
			"items":            order.OrderItems,
			"created_at":       order.CreatedAt,
		},
	})
}
