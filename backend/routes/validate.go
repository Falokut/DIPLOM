package routes

import (
	"github.com/go-playground/validator/v10"
)

type Validate struct {
	v *validator.Validate
}

func (v Validate) Validate(s any) error {
	return v.v.Struct(s)
}
