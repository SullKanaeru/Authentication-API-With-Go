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

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

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