package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model

	CustomerName    string
	CustomerPhone   string
	CustomerAddress string
	TotalPrice      int64
	Status          string
	Items           []OrderItem
}

type OrderItem struct {
	gorm.Model

	OrderID   uint
	ProductID uint
	Quantity  uint
	UnitPrice int64

	Product Product `gorm:"foreignKey:ProductID"`
}
