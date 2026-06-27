package handlers

import (
	"net/http"
	"strconv"

	"github.com/WhoIsR/bhaa_firebase_backend/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CartHandler struct {
	db *gorm.DB
}

func NewCartHandler(db *gorm.DB) *CartHandler {
	return &CartHandler{db: db}
}

type cartRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,min=1"`
}

type updateCartRequest struct {
	Quantity int `json:"quantity" binding:"required,min=1"`
}

func userIDFromContext(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}

	switch v := userID.(type) {
	case float64:
		return uint(v), true
	case uint:
		return v, true
	case int:
		return uint(v), true
	default:
		return 0, false
	}
}

func (h *CartHandler) GetCart(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "User belum login"})
		return
	}

	var items []models.CartItem
	if err := h.db.Preload("Product").Where("user_id = ?", userID).Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengambil cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"items": items,
		},
	})
}

func (h *CartHandler) AddToCart(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "User belum login"})
		return
	}

	var req cartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Data cart tidak lengkap"})
		return
	}

	var product models.Product
	if err := h.db.Where("id = ? AND is_active = ?", req.ProductID, true).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Produk tidak ditemukan"})
		return
	}
	if product.Stock < req.Quantity {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Stok produk tidak cukup"})
		return
	}

	var item models.CartItem
	err := h.db.Where("user_id = ? AND product_id = ?", userID, req.ProductID).First(&item).Error
	if err == gorm.ErrRecordNotFound {
		item = models.CartItem{UserID: userID, ProductID: req.ProductID, Quantity: req.Quantity}
		if err := h.db.Create(&item).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal menambahkan cart"})
			return
		}
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal membaca cart"})
		return
	} else {
		item.Quantity += req.Quantity
		if item.Quantity > product.Stock {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Stok produk tidak cukup"})
			return
		}
		if err := h.db.Save(&item).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengubah cart"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Produk ditambahkan ke cart"})
}

func (h *CartHandler) UpdateCartItem(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "User belum login"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID cart tidak valid"})
		return
	}

	var req updateCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Quantity tidak valid"})
		return
	}

	var item models.CartItem
	if err := h.db.Preload("Product").Where("id = ? AND user_id = ?", uint(id), userID).First(&item).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Item cart tidak ditemukan"})
		return
	}
	if req.Quantity > item.Product.Stock {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Stok produk tidak cukup"})
		return
	}

	item.Quantity = req.Quantity
	if err := h.db.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengubah cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Cart diperbarui"})
}

func (h *CartHandler) RemoveCartItem(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "User belum login"})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "ID cart tidak valid"})
		return
	}

	result := h.db.Where("id = ? AND user_id = ?", uint(id), userID).Delete(&models.CartItem{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal menghapus cart"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Item cart tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Item cart dihapus"})
}

func (h *CartHandler) ClearCart(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "User belum login"})
		return
	}

	if err := h.db.Where("user_id = ?", userID).Delete(&models.CartItem{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengosongkan cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Cart dikosongkan"})
}
