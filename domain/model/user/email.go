// Package user
package user

import (
	"strings"

	"net/mail"

	"github.com/cockroachdb/errors"
)

type Email string

// TODO: @の前後の文字数チェックをしたい。RFC基準だと複雑すぎるので
func NewEmail(val string) (Email, error) {
	if strings.TrimSpace(val) == "" {
		return "", errors.New("email is required")
	}

	parsedAddress, err := mail.ParseAddress(val)
	if err != nil {
		return "", errors.Wrap(err, "invalid email address")
	}

	email := Email(parsedAddress.Address)
	return email, nil
}

func Reconstruct(val string) (*Email, error) {
	email := Email(val)
	return &email, nil
}
