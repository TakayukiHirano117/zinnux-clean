// Package repository
package repository

import (
	"go-rest-api/domain/model/user"
	"go-rest-api/model"
	"time"

	"gorm.io/gorm"
)

type IUserRepository interface {
	GetUserByEmail(user *model.User, email string) error
	CreateUser(user *user.UserEntity, passwordHash []byte) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &userRepository{db}
}

func (ur *userRepository) GetUserByEmail(user *model.User, email string) error {
	if err := ur.db.Where("email=?", email).First(user).Error; err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) CreateUser(user *user.UserEntity, passwordHash []byte) error {
	// データモデルに詰め替える
	userDataModel := model.User{
		ID: user.id,
		Email: user.email,
		Password: string(passwordHash),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// データモデルを保存する
	if err := ur.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}
