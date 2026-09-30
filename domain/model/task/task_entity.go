// Package task
package task

import (
	"go-rest-api/domain/model/shared"
	"unicode/utf8"

	"github.com/cockroachdb/errors"
)

type TaskEntity struct {
	id     shared.UUID
	userID shared.UUID
	title  string
}

func NewTaskEntity(id shared.UUID, userID shared.UUID, title string) (*TaskEntity, error) {
	if title == "" {
		return nil, errors.New("task must not be empty")
	}

	if utf8.RuneCountInString(title) > 100 {
		return nil, errors.New("task title must be in 100 chars")
	}

	return &TaskEntity{
		id:     id,
		userID: userID,
		title:  title,
	}, nil
}

func (te *TaskEntity) ID() shared.UUID {
	return te.id
}

func (te *TaskEntity) UserID() shared.UUID {
	return te.userID
}

func (te *TaskEntity) Title() string {
	return te.title
}
