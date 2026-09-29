// Package shared
package shared

import (
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

type UUID string

func NewUUID() UUID {
	return UUID(uuid.New().String())
}

func Reconstruct(val string) (UUID, error) {
	if val == "" {
		return "", errors.New("UUID must not be empty")
	}

	if _, err := uuid.Parse(val); err != nil {
		return "", errors.New("UUID is invalid")
	}

	return UUID(val), nil
}
