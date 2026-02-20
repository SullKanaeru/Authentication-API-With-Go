package service

import (
	"authentication_api/internal/model"
	"authentication_api/internal/repository"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	Repo      *repository.UserRepository
	JWTSecret string
	Notif     *NotificationService
}

func NewAuthService(repo *repository.UserRepository, secret string, notif *NotificationService) *AuthService {
	return &AuthService{
		Repo:      repo,
		JWTSecret: secret,
		Notif:     notif,
	}
}

func generateOTP() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(1000000))
	return fmt.Sprintf("%06d", n)
}

func (s *AuthService) Register(req model.RegisterRequest) error {
	if req.Password != req.ConfirmPassword {
		return errors.New("password dan konfirmasi password tidak cocok")
	}

	existingUser, err := s.Repo.CheckExistingUser(req.Username, req.Email, req.PhoneNumber)
	
	if err == nil && existingUser != nil {
		
		if existingUser.IsVerified {
			return errors.New("pendaftaran gagal: username, email, atau nomor telepon sudah terdaftar")
		}

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		newOTP := generateOTP()

		existingUser.Password = string(hashedPassword)
		existingUser.Fullname = req.Fullname 
		existingUser.OTP = newOTP
		existingUser.OTPExpiredAt = time.Now().Add(5 * time.Minute)

		if err := s.Repo.UpdateUser(existingUser); err != nil {
			return errors.New("gagal memperbarui data registrasi")
		}

		go s.sendVerification(existingUser, newOTP, req.SendNotificationTo)
		return nil 
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("gagal memproses password")
	}

	otp := generateOTP()
	expiredAt := time.Now().Add(5 * time.Minute)

	user := &model.User{
		Fullname:     req.Fullname,
		Username:     req.Username,
		Email:        req.Email,
		PhoneNumber:  req.PhoneNumber,
		Password:     string(hashedPassword),
		Role:         "user",
		IsVerified:   false,
		OTP:          otp,
		OTPExpiredAt: expiredAt,
	}

	if err := s.Repo.CreateUser(user); err != nil {
		return errors.New("gagal membuat akun baru")
	}

	go s.sendVerification(user, otp, req.SendNotificationTo)

	return nil
}

func (s *AuthService) Login(req model.LoginRequest) (string, error) {
	user, err := s.Repo.FindByIdentifier(req.Identifier)
	if err != nil {
		return "", errors.New("kredensial tidak valid")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", errors.New("kredensial tidak valid")
	}

	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,                            
		"exp":      time.Now().Add(time.Hour * 24).Unix(), 
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.JWTSecret))
	if err != nil {
		return "", errors.New("gagal membuat token autentikasi")
	}

	return signedToken, nil
}

func (s *AuthService) sendVerification(user *model.User, otp string, sendNotificationTo string) {
	switch sendNotificationTo {
	case "wa":
		message := fmt.Sprintf("Halo %s,\n\nKode verifikasi Anda adalah: %s. Kode ini berlaku selama 5 menit. Jangan berikan kode ini kepada siapapun.", user.Fullname, otp)
		err := s.Notif.SendWhatsApp(user.PhoneNumber, message)
		if err != nil {
			log.Println("❌ Gagal kirim WA:", err)
		}

	case "email":
		subject := "Verifikasi Akun Anda"
		body := fmt.Sprintf("Halo %s,\n\nKode verifikasi Anda adalah: %s. Kode ini berlaku selama 5 menit.", user.Fullname, otp)
		err := s.Notif.SendEmail(user.Email, subject, body)
		if err != nil {
			log.Println("❌ Gagal kirim email:", err)
		}
	}
}

func (s *AuthService) VerifyOTP(req model.VerifyRequest) error {
	user, err := s.Repo.FindByIdentifier(req.Identifier)
	if err != nil {
		return errors.New("user tidak ditemukan")
	}

	if user.IsVerified {
		return errors.New("akun ini sudah diverifikasi sebelumnya")
	}

	if user.OTP != req.OTP {
		return errors.New("kode OTP salah")
	}

	if time.Now().After(user.OTPExpiredAt) {
		return errors.New("kode OTP sudah kadaluarsa, silakan minta kode baru")
	}

	user.IsVerified = true
	user.OTP = ""

	return s.Repo.UpdateUser(user)
}
