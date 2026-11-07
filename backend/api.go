package main

import (
	"fmt"
	"net/http"
	"strconv"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"sort"
	"time"

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
		api.GET("/inventory/materials", getMaterialInventory)
		api.POST("/materials/:id/stock", stockInMaterial)

		api.GET("/products", getProducts)
		api.POST("/products", createProduct)
		api.GET("/products/:id", getProduct)
		api.PUT("/products/:id", updateProduct)
		api.DELETE("/products/:id", deleteProduct)
		api.POST("/products/:id/boms", createBOM)
		api.POST("/products/:id/bom-explosion", bomExplosion)
		api.POST("/work-orders", createWorkOrder)
		api.GET("/work-orders", getWorkOrders)
		api.POST("/work-orders/:id/complete", completeWorkOrder)
		api.GET("/inventory/products", getProductInventory)
		
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

func getMaterialInventory(c *gin.Context) {
    var inventories []MaterialInventory

    if err := db.Preload("Material").Find(&inventories).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Gagal mengambil data inventory",
        })
        return
    }

    var results []gin.H

    for _, inventory := range inventories {
        available := inventory.OnHand - inventory.Reserved

        results = append(results, gin.H{
            "material_id":   inventory.MaterialID,
            "material_code": inventory.Material.Code,
            "material_name": inventory.Material.Name,
            "base_unit":     inventory.Material.BaseUnit,
            "on_hand":       inventory.OnHand,
            "reserved":      inventory.Reserved,
            "available":     available,
        })
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Data inventory berhasil diambil",
        "data":    results,
    })
}

type StockInRequest struct {
    Quantity float64 `json:"quantity" binding:"required,gt=0"`
}

func stockInMaterial(c *gin.Context) {
    materialID, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "ID material tidak valid",
        })
        return
    }

    var req StockInRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Quantity harus lebih dari 0",
        })
        return
    }

    err = db.Transaction(func(tx *gorm.DB) error {
        var inventory MaterialInventory

        // Mengunci baris inventory selama transaksi.
        if err := tx.Clauses(clause.Locking{
            Strength: "UPDATE",
        }).Where("material_id = ?", materialID).
            First(&inventory).Error; err != nil {
            return err
        }

        // Menambahkan stok bahan baku.
        inventory.OnHand += req.Quantity

        if err := tx.Save(&inventory).Error; err != nil {
            return err
        }

        return nil
    })

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Gagal menambahkan stok. Pastikan material dan inventory tersedia.",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Stok bahan baku berhasil ditambahkan",
        "data": gin.H{
            "material_id": materialID,
            "stock_added": req.Quantity,
        },
    })
}

type CreateWorkOrderRequest struct {
    ProductID uint    `json:"product_id" binding:"required"`
    Quantity  float64 `json:"quantity" binding:"required,gt=0"`
}

func createWorkOrder(c *gin.Context) {
    var req CreateWorkOrderRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Product ID dan quantity harus valid",
        })
        return
    }

    var result WorkOrder

    err := db.Transaction(func(tx *gorm.DB) error {

        // 1. Ambil BOM aktif beserta material.
        var bom BOM

        if err := tx.Preload("Product").
            Preload("Items.Material").
            Where("product_id = ? AND is_active = ?", req.ProductID, true).
            First(&bom).Error; err != nil {
            return fmt.Errorf("BOM aktif tidak ditemukan")
        }

        // 2. Hitung kebutuhan material.
        var items []WorkOrderItem

        for _, item := range bom.Items {
            if item.Unit != item.Material.BaseUnit {
                return fmt.Errorf(
                    "unit BOM tidak sesuai dengan unit material: %s",
                    item.Material.Name,
                )
            }

            required := item.Quantity * req.Quantity

            items = append(items, WorkOrderItem{
                MaterialID:       item.MaterialID,
                RequiredQuantity: required,
                IssuedQuantity:   0,
                Unit:             item.Unit,
            })
        }

        // 3. Urutkan material agar penguncian konsisten.
        sort.Slice(items, func(i, j int) bool {
            return items[i].MaterialID < items[j].MaterialID
        })

        // 4. Kunci inventory dan reservasi stok.
        for _, item := range items {
            var inventory MaterialInventory

            if err := tx.Clauses(clause.Locking{
                Strength: "UPDATE",
            }).Where("material_id = ?", item.MaterialID).
                First(&inventory).Error; err != nil {
                return fmt.Errorf(
                    "inventory material %d tidak ditemukan",
                    item.MaterialID,
                )
            }

            available := inventory.OnHand - inventory.Reserved

            if available < item.RequiredQuantity {
                return fmt.Errorf(
                    "stok material %d tidak mencukupi: tersedia %.2f, dibutuhkan %.2f",
                    item.MaterialID,
                    available,
                    item.RequiredQuantity,
                )
            }

            inventory.Reserved += item.RequiredQuantity

            if err := tx.Save(&inventory).Error; err != nil {
                return err
            }
        }

        // 5. Buat Work Order.
        result = WorkOrder{
            OrderNumber: fmt.Sprintf("WO-%d", time.Now().UnixNano()),
            ProductID:   req.ProductID,
            BOMID:       bom.ID,
            BOMVersion:  bom.Version,
            Quantity:    req.Quantity,
            Status:      "RESERVED",
        }

        if err := tx.Create(&result).Error; err != nil {
            return err
        }

        // 6. Simpan snapshot kebutuhan material.
        for i := range items {
            items[i].WorkOrderID = result.ID
        }

        if err := tx.Create(&items).Error; err != nil {
            return err
        }

        return nil
    })

    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message": "Work Order berhasil dibuat dan material berhasil direservasi",
        "data":    result,
    })
}

func getWorkOrders(c *gin.Context) {
    var workOrders []WorkOrder

    if err := db.Preload("Product").
        Preload("Items.Material").
        Order("created_at DESC").
        Find(&workOrders).Error; err != nil {

        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Gagal mengambil data Work Order",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Data Work Order berhasil diambil",
        "data":    workOrders,
    })
}

func completeWorkOrder(c *gin.Context) {
    workOrderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "ID Work Order tidak valid",
        })
        return
    }

    err = db.Transaction(func(tx *gorm.DB) error {

        // 1. Kunci Work Order.
        var workOrder WorkOrder

        if err := tx.Clauses(clause.Locking{
            Strength: "UPDATE",
        }).First(&workOrder, workOrderID).Error; err != nil {
            return fmt.Errorf("Work Order tidak ditemukan")
        }

        // 2. Pastikan Work Order belum selesai.
        if workOrder.Status != "RESERVED" {
            return fmt.Errorf("Work Order tidak dalam status RESERVED")
        }

        // 3. Ambil daftar kebutuhan material.
        var items []WorkOrderItem

        if err := tx.Where("work_order_id = ?", workOrder.ID).
            Find(&items).Error; err != nil {
            return err
        }

        if len(items) == 0 {
            return fmt.Errorf("Work Order tidak memiliki item material")
        }

        // 4. Urutkan material untuk mencegah deadlock.
        sort.Slice(items, func(i, j int) bool {
            return items[i].MaterialID < items[j].MaterialID
        })

        // 5. Potong stok bahan baku.
        for _, item := range items {
            var inventory MaterialInventory

            if err := tx.Clauses(clause.Locking{
                Strength: "UPDATE",
            }).Where("material_id = ?", item.MaterialID).
                First(&inventory).Error; err != nil {
                return err
            }

            if inventory.OnHand < item.RequiredQuantity ||
                inventory.Reserved < item.RequiredQuantity {
                return fmt.Errorf(
                    "stok material %d tidak mencukupi untuk penyelesaian",
                    item.MaterialID,
                )
            }

            inventory.OnHand -= item.RequiredQuantity
            inventory.Reserved -= item.RequiredQuantity

            if err := tx.Save(&inventory).Error; err != nil {
                return err
            }

            item.IssuedQuantity = item.RequiredQuantity

            if err := tx.Save(&item).Error; err != nil {
                return err
            }
        }

        // 6. Kunci inventory produk jadi.
        var productInventory ProductInventory

        if err := tx.Clauses(clause.Locking{
            Strength: "UPDATE",
        }).Where("product_id = ?", workOrder.ProductID).
            First(&productInventory).Error; err != nil {
            return fmt.Errorf("Inventory produk tidak ditemukan")
        }

        // 7. Tambahkan stok produk jadi.
        productInventory.OnHand += workOrder.Quantity

        if err := tx.Save(&productInventory).Error; err != nil {
            return err
        }

        // 8. Ubah status Work Order.
        workOrder.Status = "COMPLETED"

        if err := tx.Save(&workOrder).Error; err != nil {
            return err
        }

        return nil
    })

    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Work Order berhasil diselesaikan",
        "data": gin.H{
            "work_order_id": workOrderID,
            "status":        "COMPLETED",
        },
    })
}

func getProductInventory(c *gin.Context) {
    var inventories []ProductInventory

    if err := db.Preload("Product").Find(&inventories).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Gagal mengambil data inventory produk",
        })
        return
    }

    var results []gin.H

    for _, inventory := range inventories {
        available := inventory.OnHand - inventory.Reserved

        results = append(results, gin.H{
            "product_id":   inventory.ProductID,
            "product_code": inventory.Product.Code,
            "product_name": inventory.Product.Name,
            "on_hand":      inventory.OnHand,
            "reserved":     inventory.Reserved,
            "available":    available,
        })
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Data inventory produk berhasil diambil",
        "data":    results,
    })
}