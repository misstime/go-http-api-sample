// 本包用于配置 gin 内部的验证器
// 本文件包含实现 自定义验证器 的具体代码，即：请在本文件中实现自定义验证器

// see more: https://www.liwenzhou.com/posts/Go/validator_usages/

package ginvalidator

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"regexp"
)

// 自定义验证器 使用到的正则表达式
const (
	// 11位中国大陆手机号
	cnCellPhoneNumberRegexpString = `^(13[0-9]|` + // 130-139
		`14[57]|` + // 145,147
		`15[0-35-9]|` + // 150-153,155-159
		`17[0678]|` + // 170,176,177,178
		`18[0-9])` + // 180-189
		`[0-9]{8}$`
)

var (
	cnCellPhoneNumberRegexp = regexp.MustCompile(cnCellPhoneNumberRegexpString)
)

// validationArguments 类型定义：注册自定义验证器时所需的参数
type validationArguments = struct {
	tag       string
	fn        validator.Func
	errorText string
}

// batchValidationArguments 定义所有自定义验证器参数
//
// 第一个参数：验证器名称
// 第二个参数：验证器校验函数
// 第三个参数：验证器中文翻译模板
//
// Note: 此变量为实现自定义验证器的入口，****所有自定义验证器都应当在此注册****。
// 注册自定义验证器的具体实现，请参见 `cnCellPhoneNumber` 实例
var batchValidationArguments = []validationArguments{
	{"cnCellPhoneNumber", isCnCellPhoneNumber, "{0}必须是有效的11位中国大陆手机号"},
}

// isCnCellPhoneNumber 是自定义验证器 `cnCellPhoneNumber` 的校验函数：11位中国大陆手机号
var isCnCellPhoneNumber validator.Func = func(fl validator.FieldLevel) bool {
	return isMatch(cnCellPhoneNumberRegexp, fl.Field().Interface())
}

// isMatch 是一个通用的正则验证助手函数（支持类型：string、数字、[]byte、[]rune）
func isMatch(reg *regexp.Regexp, v interface{}) bool {
	switch v.(type) {
	case string:
		return reg.MatchString(v.(string))
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return reg.MatchString(fmt.Sprintf("%d", v))
	case []byte:
		return reg.Match(v.([]byte))
	case []rune:
		return reg.MatchString(string(v.([]rune)))
	default:
		return false
	}
}
