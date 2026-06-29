package handlers

import (
	"net/http"
	"strconv"

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

// GetMyOrders — ambil daftar riwayat pesanan user login (dengan nama produk via batch query)
func (h *OrderHandler) GetMyOrders(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "User belum login"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 10
	}

	orders, err := h.service.GetOrdersByUserID(userID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengambil riwayat order"})
		return
	}

	if len(orders) == 0 {
		c.JSON(http.StatusOK, gin.H{"success": true, "data": []any{}})
		return
	}

	// Batch load product names untuk semua item di semua order
	seen := make(map[uint]string)
	var ids []uint
	for _, order := range orders {
		for _, item := range order.OrderItems {
			if _, ok := seen[item.ProductID]; !ok {
				seen[item.ProductID] = ""
				ids = append(ids, item.ProductID)
			}
		}
	}
	if len(ids) > 0 {
		var products []struct {
			ID   uint
			Name string
		}
		config.DB.Model(&models.Product{}).Select("id, name").Where("id IN ?", ids).Find(&products)
		for _, p := range products {
			seen[p.ID] = p.Name
		}
	}

	type itemJSON struct {
		ProductID  uint    `json:"product_id"`
		ProductName string `json:"product_name"`
		Price      float64 `json:"price"`
		Quantity   int     `json:"quantity"`
		Subtotal   float64 `json:"subtotal"`
	}

	type orderJSON struct {
		ID              uint       `json:"id"`
		TotalAmount     float64    `json:"total_amount"`
		Status          string     `json:"status"`
		ShippingAddress string     `json:"shipping_address"`
		Notes           string     `json:"notes"`
		PaymentMethod   string     `json:"payment_method"`
		Items           []itemJSON `json:"items"`
		CreatedAt       string     `json:"created_at"`
	}

	result := make([]orderJSON, 0, len(orders))
	for _, order := range orders {
		items := make([]itemJSON, 0, len(order.OrderItems))
		for _, item := range order.OrderItems {
			items = append(items, itemJSON{
				ProductID:   item.ProductID,
				ProductName: seen[item.ProductID],
				Price:       item.Price,
				Quantity:    item.Quantity,
				Subtotal:    float64(item.Quantity) * item.Price,
			})
		}

		createdAt := ""
		if !order.CreatedAt.IsZero() {
			createdAt = order.CreatedAt.Format("2006-01-02T15:04:05Z07:00")
		}

		result = append(result, orderJSON{
			ID:          order.ID,
			TotalAmount: order.TotalPrice,
			Status:      order.Status,
			Items:       items,
			CreatedAt:   createdAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}
