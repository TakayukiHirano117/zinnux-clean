// Package user
package user

import (
	"go-rest-api/domain/model/shared"
	"net/mail"
	"strings"

	"github.com/cockroachdb/errors"
)

type UserEntity struct {
	id    shared.UUID
	email string
}

func NewUserEntity(id shared.UUID, email string) (*UserEntity, error) {
	// TODO: これはVOにするので後で消す
	trimmedEmail := strings.TrimSpace(email)
	if trimmedEmail == "" {
		return nil, errors.New("email is required")
	}

	// TODO: これはVOにするので後で消す
	parsedEmail, err := mail.ParseAddress(trimmedEmail)
	if err != nil {
		return nil, errors.Wrap(err, "invalid email address")
	}

	return &UserEntity{
		id:    id,
		email: parsedEmail.Address,
	}, nil
}

func (ue *UserEntity) ID() shared.UUID {
	return ue.id
}

func (ue *UserEntity) Email() string {
	return ue.email
}
