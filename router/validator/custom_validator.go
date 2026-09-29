// Package validator
package validator

import "github.com/go-playground/validator/v10"

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		// エラーを全て構造体に入れて返すとかしたいが。
		return err
	}
	return nil
}
