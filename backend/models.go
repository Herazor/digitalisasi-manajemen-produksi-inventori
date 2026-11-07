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