// Package user
package user

import (
	"go-rest-api/domain/model/shared"
)

type UserEntity struct {
	id    shared.UUID
	email Email
}

// const maxUnCompletedTaskCount = 20

func NewUserEntity(id shared.UUID, email Email) (*UserEntity, error) {
	return &UserEntity{
		id:    id,
		email: email,
	}, nil
}

// これ現状のままだとドメインサービスの方が適切かもしれん。ユーザーとタスク別集約なので。
// func (ue *UserEntity) CanCreateTask(uncompletedTaskCount int) bool {
// 	canCreate :=  uncompletedTaskCount < maxUnCompletedTaskCount
// 	return canCreate
// }

func (ue *UserEntity) ID() shared.UUID {
	return ue.id
}

func (ue *UserEntity) Email() Email {
	return ue.email
}
