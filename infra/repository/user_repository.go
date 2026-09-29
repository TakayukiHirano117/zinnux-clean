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

func (ur *userRepository) CreateUser(ue *user.UserEntity, passwordHash []byte) error {
	// データモデルに詰め替える
	userDataModel := model.User{
		ID: ue.ID(),
		Email: ue.Email(),
		Password: string(passwordHash),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := ur.db.Create(userDataModel).Error; err != nil {
		return err
	}
	return nil
}
