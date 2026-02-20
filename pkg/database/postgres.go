package database

import (
	"log"
	"authentication_api/internal/config"
	"authentication_api/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {
	dsn := config.GetEnv("DB_URL", "host=localhost user=postgres password=rahasia dbname=mydb port=5432 sslmode=disable TimeZone=Asia/Jakarta")

	// 2. Buka Koneksi dengan GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal koneksi ke database:", err)
	}

	log.Println("Database (Supabase) terhubung!")
	return db
}

func MigrateDB(db *gorm.DB) {
	err := db.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatal("Gagal melakukan migrasi:", err)
	}
	log.Println("Migrasi database berhasil!")
}