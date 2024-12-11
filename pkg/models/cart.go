package models

import "gorm.io/gorm"

type CartItem struct {
	gorm.Model
	//ID          uint    `gorm:"primaryKey"`
	UserID      string  `gorm:"not null"`
	ProductID   string  `gorm:"not null"`
	ProductName string  `gorm:"not null"`
	Price       float64 `gorm:"not null"`
	Quantity    int32   `gorm:"not null"`
}

type Product struct {
	ID           uint    `gorm:"primaryKey"`
	ProductName  string  `gorm:"not null"`
	Description  string  `gorm:"not null"`
	ImageUrl     string  `gorm:"not null"`
	Price        float64 `gorm:"not null"`
	Stock        int32   `gorm:"not null"`
	CategoryName string  `gorm:"not null"`
}
