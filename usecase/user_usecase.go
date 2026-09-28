// Package usecase
package usecase

import (
	"go-rest-api/domain/model/user"
	"go-rest-api/model"
	"go-rest-api/repository"
	"os"
	"time"
	"github.com/google/uuid"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type IUserUsecase interface {
	SignUp(signUpRequestDTO SignUpRequestDTO) (SignUpResponseDTO, error)
	Login(user model.User) (string, error)
}

type userUsecase struct {
	ur repository.IUserRepository
}

func NewUserUsecase(ur repository.IUserRepository) IUserUsecase {
	return &userUsecase{ur}
}

type SignUpRequestDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpResponseDTO struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

func (uu *userUsecase) SignUp(signUpRequestDTO SignUpRequestDTO) (model.UserResponse, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(signUpRequestDTO.Password), 10)
	if err != nil {
		return model.UserResponse{}, err
	}
	// newUser := model.User{Email: user.Email, Password: string(hash)}
	newUser, err := user.NewUserEntity(uuid.New(), signUpRequestDTO.Email)
	if err := uu.ur.CreateUser(&newUser); err != nil {
		// errors出す
		return model.UserResponse{}, err
	}

	resUser := model.UserResponse{
		ID:    newUser.ID,
		Email: newUser.Email,
	}
	return resUser, nil
}

func (uu *userUsecase) Login(user model.User) (string, error) {
	storedUser := model.User{}
	if err := uu.ur.GetUserByEmail(&storedUser, user.Email); err != nil {
		return "", err
	}

	err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": storedUser.ID,
		"exp":     time.Now().Add(time.Hour * 12).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return "", nil
	}
	return tokenString, nil
}
