package user

import (
	"github.com/changaolee/skeleton/pkg/validation"
	"github.com/changaolee/skeleton/pkg/validation/field"
)

// Validate 检查一个 user 对象是否合法.
func (u *User) Validate() field.ErrorList {
	val := validation.NewValidator(u)
	allErrs := val.Validate()

	if err := validation.IsValidPassword(u.Password); err != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("password"), err.Error(), ""))
	}

	return allErrs
}
