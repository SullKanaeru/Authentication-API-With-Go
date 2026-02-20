package model

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Fullname     string    `gorm:"type:varchar(255);not null" json:"fullname"`
	Username     string    `gorm:"type:varchar(100);not null" json:"username"`
	Email        string    `gorm:"type:varchar(150);not null" json:"email"`
	PhoneNumber  string    `gorm:"type:varchar(20);not null" json:"phone_number"`
	Password     string    `gorm:"type:varchar(100);not null" json:"password"`
	Role         string    `gorm:"type:varchar(20);not null" json:"role"`
	IsVerified   bool      `gorm:"default:false" json:"is_verified"`
	OTP          string    `gorm:"type:varchar(10)" json:"-"`
	OTPExpiredAt time.Time `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type RegisterRequest struct {
	Fullname        string `json:"fullname" binding:"required"`
	Username        string `json:"username" binding:"required"`
	Email           string `json:"email" binding:"required"`
	PhoneNumber     string `json:"phone_number" binding:"required"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
	SendNotificationTo string `json:"send_notification_to" binding:"required"`
}

type LoginRequest struct {
	Identifier string `json:"identifier" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

type VerifyRequest struct {
	Identifier string `json:"identifier" binding:"required"`
	OTP        string `json:"otp" binding:"required"`
}
