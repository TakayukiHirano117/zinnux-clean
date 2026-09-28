package shared

import "golang.org/x/crypto/bcrypt"

type IPasswordGenerator interface {
	Execute(plainTextPassword string, cost int) ([]byte, error)
}

type passwordGenerator struct{}


func NewPasswordGenerator() IPasswordGenerator {
	return &passwordGenerator{}
}

func (pg *passwordGenerator) Execute(plainTextPassword string, cost int) ([]byte, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), cost)
	if err != nil {
		return nil, err
	}
	return passwordHash, nil
}
