package models

import "gorm.io/gorm"

type CartItem struct {
	gorm.Model
	UserID    uint    `gorm:"not null;index:idx_user_product,unique" json:"user_id"`
	ProductID uint    `gorm:"not null;index:idx_user_product,unique" json:"product_id"`
	Product   Product `gorm:"foreignKey:ProductID" json:"product"`
	Quantity  int     `gorm:"not null;default:1" json:"quantity"`
}
