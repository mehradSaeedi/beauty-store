package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model

	CustomerName    string
	CustomerPhone   string
	CustomerAddress string
	TotalPrice      int64
	Items           []OrderItem
}

type OrderItem struct {
	gorm.Model

	OrderID   uint
	ProductID uint
	Quantity  uint
	UnitPrice int64
}
