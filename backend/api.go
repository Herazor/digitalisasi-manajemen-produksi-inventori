package main

import (
	"fmt"
	"net/http"
	"strconv"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

var db *gorm.DB

func setupRoutes() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		api.GET("/materials", getMaterials)
		api.POST("/materials", createMaterial)
		api.GET("/materials/:id", getMaterial)
		api.PUT("/materials/:id", updateMaterial)
		api.DELETE("/materials/:id", deleteMaterial)

		api.GET("/products", getProducts)
		api.POST("/products", createProduct)
		api.GET("/products/:id", getProduct)
		api.PUT("/products/:id", updateProduct)
		api.DELETE("/products/:id", deleteProduct)
		api.POST("/products/:id/boms", createBOM)
		api.POST("/products/:id/bom-explosion", bomExplosion)
	}

	return r
}

func getMaterials(c *gin.Context) {
	var materials []Material

	if err := db.Find(&materials).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": materials})
}

func createMaterial(c *gin.Context) {
	var material Material

	if err := c.ShouldBindJSON(&material); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}

	if material.Code == "" || material.Name == "" || material.BaseUnit == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Code, name, dan base_unit wajib diisi"})
		return
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&material).Error; err != nil {
			return err
		}

		inventory := MaterialInventory{MaterialID: material.ID}
		return tx.Create(&inventory).Error
	})

	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Gagal membuat material: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Material berhasil dibuat", "data": material})
}

func getMaterial(c *gin.Context) {
	id := c.Param("id")
	var material Material

	if err := db.First(&material, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Material tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": material})
}

func updateMaterial(c *gin.Context) {
	id := c.Param("id")
	var material Material

	if err := db.First(&material, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Material tidak ditemukan"})
		return
	}

	var input struct {
		Code        string `json:"code"`
		Name        string `json:"name"`
		BaseUnit    string `json:"base_unit"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}

	material.Code = input.Code
	material.Name = input.Name
	material.BaseUnit = input.BaseUnit
	material.Description = input.Description

	if err := db.Save(&material).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Material berhasil diperbarui", "data": material})
}

func deleteMaterial(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	if err := db.Delete(&Material{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Material berhasil dihapus"})
}

func getProducts(c *gin.Context) {
	var products []Product

	if err := db.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": products})
}

func createProduct(c *gin.Context) {
	var product Product

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}

	if product.Code == "" || product.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Code dan name wajib diisi"})
		return
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&product).Error; err != nil {
			return err
		}

		inventory := ProductInventory{ProductID: product.ID}
		return tx.Create(&inventory).Error
	})

	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Gagal membuat product: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Product berhasil dibuat", "data": product})
}

func getProduct(c *gin.Context) {
	id := c.Param("id")
	var product Product

	if err := db.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": product})
}

func updateProduct(c *gin.Context) {
	id := c.Param("id")
	var product Product

	if err := db.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product tidak ditemukan"})
		return
	}

	var input struct {
		Code        string `json:"code"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak valid"})
		return
	}

	product.Code = input.Code
	product.Name = input.Name
	product.Description = input.Description

	if err := db.Save(&product).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product berhasil diperbarui", "data": product})
}

func deleteProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	if err := db.Delete(&Product{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product berhasil dihapus"})
}

type CreateBOMRequest struct {
	Items []struct {
		MaterialID uint    `json:"material_id"`
		Quantity   float64 `json:"quantity"`
		Unit       string  `json:"unit"`
	} `json:"items"`
}

func createBOM(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID product tidak valid",
		})
		return
	}

	var req CreateBOMRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Format JSON tidak valid",
		})
		return
	}

	if len(req.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "BOM harus memiliki minimal satu bahan baku",
		})
		return
	}

	var createdBOM BOM

	err = db.Transaction(func(tx *gorm.DB) error {
		var product Product

		if err := tx.First(&product, uint(productID)).Error; err != nil {
			return err
		}

		var bomCount int64

if err := tx.Model(&BOM{}).
    Where("product_id = ?", productID).
    Count(&bomCount).Error; err != nil {
    return err
}

version := 1

if bomCount > 0 {
    var lastBOM BOM

    if err := tx.Where("product_id = ?", productID).
        Order("version DESC").
        First(&lastBOM).Error; err != nil {
        return err
    }

    version = lastBOM.Version + 1
}

		for _, item := range req.Items {
			if item.MaterialID == 0 || item.Quantity <= 0 || item.Unit == "" {
				return fmt.Errorf("data bahan baku tidak valid")
			}

			var material Material

			if err := tx.First(&material, item.MaterialID).Error; err != nil {
				return err
			}
		}

		if err := tx.Model(&BOM{}).
			Where("product_id = ?", productID).
			Update("is_active", false).Error; err != nil {
			return err
		}

		createdBOM = BOM{
			ProductID: uint(productID),
			Version:   version,
			IsActive:  true,
		}

		if err := tx.Create(&createdBOM).Error; err != nil {
			return err
		}

		for _, item := range req.Items {
			bomItem := BOMItem{
				BOMID:      createdBOM.ID,
				MaterialID: item.MaterialID,
				Quantity:   item.Quantity,
				Unit:       item.Unit,
			}

			if err := tx.Create(&bomItem).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	db.Preload("Items.Material").
		Preload("Product").
		First(&createdBOM, createdBOM.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "BOM berhasil dibuat",
		"data":    createdBOM,
	})
}

type BOMExplosionRequest struct {
	Quantity float64 `json:"quantity" binding:"required,gt=0"`
}

func bomExplosion(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID product tidak valid",
		})
		return
	}

	var req BOMExplosionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Quantity harus lebih dari 0",
		})
		return
	}

	var bom BOM

	err = db.Preload("Product").
		Preload("Items.Material").
		Where("product_id = ? AND is_active = ?", productID, true).
		First(&bom).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "BOM aktif tidak ditemukan",
		})
		return
	}

	var requirements []gin.H

	for _, item := range bom.Items {
		totalQuantity := item.Quantity * req.Quantity

		requirements = append(requirements, gin.H{
			"material_id":   item.MaterialID,
			"material_code": item.Material.Code,
			"material_name": item.Material.Name,
			"quantity_per_product": item.Quantity,
			"total_quantity": totalQuantity,
			"unit": item.Unit,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "BOM explosion berhasil",
		"data": gin.H{
			"product_id":   bom.ProductID,
			"product_name": bom.Product.Name,
			"bom_version":  bom.Version,
			"production_quantity": req.Quantity,
			"materials": requirements,
		},
	})
}