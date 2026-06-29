package models

import "gorm.io/gorm"

type CartItem struct {
	gorm.Model
	UserID    uint    `gorm:"not null" json:"user_id"`
	ProductID uint    `gorm:"not null" json:"product_id"`
	Product   Product `gorm:"foreignKey:ProductID" json:"product"`
	Quantity  int     `gorm:"not null;default:1" json:"quantity"`
}
