// Package user
package user

import (
	"github.com/google/uuid"
)

type UserEntity struct {
	id    uuid.UUID
	email string
}

func NewUserEntity(id uuid.UUID, email string) (*UserEntity, error) {
	// UUIDかどうか
	// メールアドレスのバリデーション
	return &UserEntity{
		id:    id,
		email: email,
	}, nil
}

// getter定義
func (ue *UserEntity) ID() uuid.UUID {
	return ue.id
}
