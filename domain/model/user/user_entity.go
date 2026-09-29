// Package user
package user

import (
	"go-rest-api/domain/model/shared"
)

type UserEntity struct {
	id    shared.UUID
	email Email
}

func NewUserEntity(id shared.UUID, email Email) (*UserEntity, error) {
	return &UserEntity{
		id:    id,
		email: email,
	}, nil
}

func (ue *UserEntity) ID() shared.UUID {
	return ue.id
}

func (ue *UserEntity) Email() Email {
	return ue.email
}
