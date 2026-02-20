package repository

import (
	"errors"
	"authentication_api/internal/model"
	"strings"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) CreateUser(user *model.User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepository) FindByIdentifier(identifier string) (*model.User, error) {
	var user model.User
	
	query := "phone_number = ?"
	if strings.Contains(identifier, "@") {
		query = "email = ?"
	}

	err := r.DB.Where(query, identifier).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user tidak ditemukan")
		}
		return nil, err
	}
	
	return &user, nil
}

func (r *UserRepository) UpdateUser(user *model.User) error {
	return r.DB.Save(user).Error
}

func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}