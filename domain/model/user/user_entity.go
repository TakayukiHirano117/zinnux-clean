// Package user
package user

import (
	"net/mail"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

type UserEntity struct {
	id    uuid.UUID
	email string
}

func NewUserEntity(id uuid.UUID, email string) (*UserEntity, error) {
	if id == uuid.Nil {
		return nil, errors.New("user id must be a valid UUID")
	}

	// TODO: これはVOにするので後で消す
	trimmedEmail := strings.TrimSpace(email)
	if trimmedEmail == "" {
		return nil, errors.New("email is required")
	}

	parsedEmail, err := mail.ParseAddress(trimmedEmail)
	if err != nil {
		return nil, errors.Wrap(err, "invalid email address")
	}

	return &UserEntity{
		id:    id,
		email: parsedEmail.Address,
	}, nil
}

func (ue *UserEntity) ID() uuid.UUID {
	return ue.id
}

func (ue *UserEntity) Email() string {
	return ue.email
}
