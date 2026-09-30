// Package usecase
package usecase

import (
	domain_model_shared "go-rest-api/domain/model/shared"
	"go-rest-api/domain/model/user"
	"go-rest-api/infra/repository"
	infra_shared "go-rest-api/infra/shared"
	"go-rest-api/infra/model"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type IUserUsecase interface {
	SignUp(signUpRequestDTO *SignUpRequestDTO) (*SignUpResponseDTO, error)
	Login(user model.User) (string, error)
}

type userUsecase struct {
	ur repository.IUserRepository
	pg infra_shared.IPasswordGenerator
}

func NewUserUsecase(
	ur repository.IUserRepository,
	pg infra_shared.IPasswordGenerator,
) IUserUsecase {
	return &userUsecase{
		ur: ur,
		pg: pg,
	}
}

type SignUpRequestDTO struct {
	Email    string
	Password string
}

const DEFAULT_COST = 10

type SignUpResponseDTO struct {
	ID domain_model_shared.UUID
}

func (uu *userUsecase) SignUp(signUpRequestDTO *SignUpRequestDTO) (*SignUpResponseDTO, error) {
	email, err := user.NewEmail(signUpRequestDTO.Email)
	newUser, err := user.NewUserEntity(domain_model_shared.NewUUID(), email)
	if err != nil {
		return nil, err
	}

	passwordHash, err := uu.pg.Execute(signUpRequestDTO.Password, DEFAULT_COST)
	if err != nil {
		return nil, err
	}

	if err := uu.ur.CreateUser(newUser, passwordHash); err != nil {
		return nil, err
	}

	return &SignUpResponseDTO{ID: newUser.ID()}, nil
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
