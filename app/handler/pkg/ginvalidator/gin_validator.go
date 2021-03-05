// 本包用于配置 gin 内部的验证器

// see more: https://www.liwenzhou.com/posts/Go/validator_usages/

package ginvalidator

import (
	"encoding/json"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"strings"
)

// Init 用于初始化 gin 内部的验证器
func Init() error {
	v := binding.Validator.Engine().(*validator.Validate)
	if err := initTrans(v, "zh"); err != nil {
		return err
	}
	initTagName(v)
	if err := registerValidations(v); err != nil {
		return err
	}
	return nil
}

// ValidationError 类型表示最终返回的字段错误信息，是对 validator.ValidationErrors 的再加工
//
// 例：map[sting]string{"email": "email为必填字段", "age":"age必须大于或等于1"}
type ValidationError map[string]string

// Replace 用于替换 ValidationError 的字段错误信息为自定义信息。
//
// 例：校验如下结构体：
//
// type SignUpParam struct {
//	Password        string `json:"password" binding:"required"`
//	RePassword      string `json:"re_password" binding:"required,eqfield=Password"`
// }
//
// 如果入参 re_password 和 password 不一致，会得到错误如下：
// {"re_password":"re_password必须等于Password"}} -- 注意此处`Password`首字母为大写，与预期不符。
//
// 此时可以使用使用 Replace(map[string]string{"re_password":"re_password必须等于password"}) 进行替换。
// 最终 ValidationError 为： {"re_password":"re_password必须等于password"}} -- `password` 首字母小写
//
// 注：Replace() 仅是解决上述问题的一种比较简单的方法，也可以通过“自定义结构体校验方法”来解决，
// see more：https://www.liwenzhou.com/posts/Go/validator_usages/#autoid-1-0-3
func (errs ValidationError) Replace(info map[string]string) error {
	for k, v := range info {
		if _, ok := errs[k]; !ok {
			return errors.Errorf("ValidationError 错误消息替换失败：不存在 key - `%s`", k)
		} else {
			errs[k] = v
		}
	}
	return nil
}

func (errs ValidationError) Error() string {
	if byteSlice, err := json.Marshal(errs); err != nil {
		return ""
	} else {
		return string(byteSlice)
	}
}

// trimStructName 去除字段错误信息中 key 内包含的结构体名称
//
// 例："SignUpParam.password":"password为必填字段" -> "password":"password为必填字段"
//
// See more: https://www.liwenzhou.com/posts/Go/validator_usages/#autoid-1-0-2
func trimStructName(fields map[string]string) ValidationError {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}
	return res
}

// Translate 翻译一个 validator.ValidationErrors，并去除字段错误信息中 key 内包含的结构体名称。
func Translate(errs *validator.ValidationErrors) ValidationError {
	return trimStructName(errs.Translate(trans))
}
