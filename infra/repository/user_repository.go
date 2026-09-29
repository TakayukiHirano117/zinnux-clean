// Package repository
package repository

import (
	"go-rest-api/domain/model/user"
	"go-rest-api/infra/model"
	"time"

	"github.com/google/uuid"
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
	parsedUserID, err := uuid.Parse(string(ue.ID()))
	if err != nil {
		return err
	}

	userDataModel := model.User{
		ID:        parsedUserID,
		Email:     string(ue.Email()),
		Password:  string(passwordHash),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := ur.db.Create(userDataModel).Error; err != nil {
		return err
	}
	return nil
}
