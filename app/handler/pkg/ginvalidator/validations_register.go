// 本包用于配置 gin 内部的验证器

// see more: https://www.liwenzhou.com/posts/Go/validator_usages/

package ginvalidator

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

// registerValidations 注册所有自定义验证器
func registerValidations(v *validator.Validate) error {
	var err error
	for _, args := range batchValidationArguments {
		err = registerValidation(v, args.tag, args.fn, args.errorText)
		if err != nil {
			return err
		}
	}
	return nil
}

// registerValidation 注册单个自定义验证器
func registerValidation(v *validator.Validate, tag string, fn validator.Func, errorText string) error {
	if err := v.RegisterValidation(tag, fn); err != nil {
		return errors.Wrapf(err, "register validation '%s' failed", tag)
	}
	err := v.RegisterTranslation(
		tag,
		trans,
		func(ut ut.Translator) error {
			return ut.Add(tag, errorText, true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T(tag, fe.Field())
			return t
		},
	)
	if err != nil {
		return errors.Wrapf(err, "register translation '%s' failed", tag)
	}
	return nil
}
