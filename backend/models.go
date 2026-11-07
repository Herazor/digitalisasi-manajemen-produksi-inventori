package main

import (
	"time"

	"gorm.io/gorm"
)

type Material struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Code        string         `gorm:"uniqueIndex;not null" json:"code"`
	Name        string         `gorm:"not null" json:"name"`
	BaseUnit    string         `gorm:"not null" json:"base_unit"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type Product struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Code        string         `gorm:"uniqueIndex;not null" json:"code"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type MaterialInventory struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	MaterialID uint      `gorm:"uniqueIndex;not null" json:"material_id"`
	OnHand     float64   `gorm:"not null;default:0" json:"on_hand"`
	Reserved   float64   `gorm:"not null;default:0" json:"reserved"`
	UpdatedAt  time.Time `json:"updated_at"`

	Material Material `gorm:"foreignKey:MaterialID" json:"material,omitempty"`
}

type ProductInventory struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProductID uint      `gorm:"uniqueIndex;not null" json:"product_id"`
	OnHand    float64   `gorm:"not null;default:0" json:"on_hand"`
	Reserved  float64   `gorm:"not null;default:0" json:"reserved"`
	UpdatedAt time.Time `json:"updated_at"`

	Product Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

type BOM struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    ProductID uint      `gorm:"not null;uniqueIndex:idx_product_version,priority:1" json:"product_id"`
    Version   int       `gorm:"not null;uniqueIndex:idx_product_version,priority:2" json:"version"`
    IsActive  bool      `gorm:"default:true" json:"is_active"`
    CreatedAt time.Time `json:"created_at"`

    Product Product   `gorm:"foreignKey:ProductID" json:"product,omitempty"`
    Items   []BOMItem `gorm:"foreignKey:BOMID" json:"items,omitempty"`
}

type BOMItem struct {
    ID         uint    `gorm:"primaryKey" json:"id"`
    BOMID      uint    `gorm:"not null" json:"bom_id"`
    MaterialID uint    `gorm:"not null" json:"material_id"`
    Quantity   float64 `gorm:"not null" json:"quantity"`
    Unit       string  `gorm:"not null" json:"unit"`

    Material Material `gorm:"foreignKey:MaterialID" json:"material,omitempty"`
}

type WorkOrder struct {
    ID          uint   `gorm:"primaryKey" json:"id"`
    OrderNumber string `gorm:"uniqueIndex;not null" json:"order_number"`
    ProductID   uint   `json:"product_id"`
    BOMID       uint   `json:"bom_id"`
    BOMVersion  int    `json:"bom_version"`
    Quantity    float64 `json:"quantity"`
    Status      string `json:"status"`

    Product Product `gorm:"foreignKey:ProductID" json:"product"`
    Items   []WorkOrderItem `gorm:"foreignKey:WorkOrderID" json:"items"`

    CreatedAt time.Time `json:"created_at"`
}

type WorkOrderItem struct {
    ID                uint    `gorm:"primaryKey" json:"id"`
    WorkOrderID       uint    `json:"work_order_id"`
    MaterialID        uint    `json:"material_id"`
    RequiredQuantity  float64 `json:"required_quantity"`
    IssuedQuantity    float64 `json:"issued_quantity"`
    Unit              string  `json:"unit"`

    Material Material `gorm:"foreignKey:MaterialID" json:"material"`
}

// Sales Order
type SalesOrder struct {
    ID           uint      `gorm:"primaryKey" json:"id"`
    OrderNumber  string    `gorm:"uniqueIndex;not null" json:"order_number"`
    CustomerName string    `gorm:"not null" json:"customer_name"`
    Status       string    `gorm:"not null;default:RESERVED" json:"status"`

    Items []SalesOrderItem `gorm:"foreignKey:SalesOrderID" json:"items"`

    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// Sales Order Item
type SalesOrderItem struct {
    ID           uint    `gorm:"primaryKey" json:"id"`
    SalesOrderID uint    `gorm:"not null;index" json:"sales_order_id"`
    ProductID    uint    `gorm:"not null" json:"product_id"`
    Quantity     float64 `gorm:"not null" json:"quantity"`

    Product Product `gorm:"foreignKey:ProductID" json:"product"`
}