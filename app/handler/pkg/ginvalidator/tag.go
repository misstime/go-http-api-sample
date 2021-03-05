// 本包用于配置 gin 内部的验证器

// see more: https://www.liwenzhou.com/posts/Go/validator_usages/

package ginvalidator

import (
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

// initTagName 自定义错误提示信息的字段名：设置错误提示中的字段 tag 为结构体注解中的 json tag
//
// 例：校验如下结构体：
//
// type SignUpParam struct {
//	Email      string `json:"email" binding:"required,email"`
// }
//
// 不调用 initTagName，验证失败得到如下错误信息：
// {"SignUpParam.Email":"Email必须是一个有效的邮箱"} -- 注意第一个`Email`首字母大写
// 调用 initTagName，验证失败得到如下错误信息：
// {"SignUpParam.email":"Email必须是一个有效的邮箱"} -- 注意第一个`Email`首字母小写
//
// See more: https://www.liwenzhou.com/posts/Go/validator_usages/#autoid-1-0-2
func initTagName(v *validator.Validate) {
	// 注册一个获取json tag的自定义方法
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}
