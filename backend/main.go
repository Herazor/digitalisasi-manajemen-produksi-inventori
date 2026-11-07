package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)

	var err error

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal menghubungkan database:", err)
	}

	fmt.Println("PostgreSQL berhasil terhubung!")

	err = db.AutoMigrate(
    	&Material{},
    	&Product{},
   		&MaterialInventory{},
   		&ProductInventory{},
    	&BOM{},
    	&BOMItem{},
		&WorkOrder{},
    	&WorkOrderItem{},
		&SalesOrder{},
    	&SalesOrderItem{},
)

	if err != nil {
		log.Fatal("Gagal membuat tabel:", err)
	}

	fmt.Println("Database migration berhasil!")

	r := setupRoutes()

	fmt.Println("Backend berjalan di http://localhost:9000")
	log.Fatal(r.Run(":9000"))
}