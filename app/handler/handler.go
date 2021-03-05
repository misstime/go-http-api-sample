// 本包定义各种 gin.HandlerFunc
// 本文件定义一些基础函数

package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"
	"project/app/handler/pkg/e"
	"project/app/handler/pkg/ginvalidator"
)

// body 即 response body
type body struct {
	Code    e.Code           `json:"code"`
	Status  string           `json:"status"`
	Message string           `json:"message"`
	Data    interface{}      `json:"data,omitempty"`
	Error   []e.IErrorDetail `json:"error,omitempty"`
}

const (
	// gin.Context 中 body 对应的 key，供日志中间件使用
	contextKeyBody = "body"
	// gin.Context 中“错误信息”对应的 key，供日志中间件使用
	contextKeyError = "error"
	// gin.Context 中“日志级别”对应的 key，供日志中间件使用
	contextKeyLogLevel = "logLevel"
)

// logError 将请求过程中的 “错误信息” 和 “日志级别” 附到 gin.Context，供日志中间件使用。
//
// 该方法一般用于 success()、fail() 方法的间接调用。
func logError(c *gin.Context, err error, logLevel zapcore.Level) {
	c.Set(contextKeyError, err)
	c.Set(contextKeyLogLevel, logLevel)
}

// success 响应成功，无错误信息
func success(c *gin.Context, data interface{}) {
	codeDetail := e.GetCodeDetail(e.CodeOK)
	body := &body{
		Code:    codeDetail.Code,
		Status:  codeDetail.Status,
		Message: codeDetail.Message,
		Data:    data,
		Error:   nil,
	}
	c.Set(contextKeyBody, body)
	logError(c, nil, codeDetail.LogLevel)
	c.JSON(codeDetail.HttpStatus, body)
}

// fail 响应错误
func fail(c *gin.Context, err error, code e.Code, errorDetails ...e.IErrorDetail) {
	codeDetail := e.GetCodeDetail(code)
	body := &body{
		Code:    codeDetail.Code,
		Status:  codeDetail.Status,
		Message: codeDetail.Message,
		Data:    nil,
		Error:   errorDetails,
	}
	c.Set(contextKeyBody, body)
	logError(c, err, codeDetail.LogLevel)
	c.JSON(codeDetail.HttpStatus, body)
}

// mustBind 类同 c.Bind()，将 request 参数绑定到指定结构体，并进行校验。
//
// 绑定校验成功：返回 true
// 绑定校验过程中程序出错：响应 e.CodeInternal 错误
// 绑定校验过程中发现 validator.ValidationErrors 错误，响应 e.CodeInvalidArgument 错误
//
// 参数 replaces 用于自定义字段错误消息，具体用法见：project/app/handler/pkg/ginvalidator
// 包 ValidationError.Replace() 方法。
// 参数 replaces 仅第一个元素有效，且第一个元素不能为 nil
func mustBind(c *gin.Context, obj interface{}, replaces ...map[string]string) (success bool) {
	var replace map[string]string
	if replaces != nil {
		if replaces[0] != nil {
			replace = replaces[0]
		} else {
			err := errors.New("ctrl mustBind(): " +
				"参数 replaces 只有第一个元素有效，且第一个元素不能为 nil")
			fail(c, err, e.CodeInternal)
			return false
		}
	}

	err := c.ShouldBind(obj)
	if err == nil {
		return true
	}

	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		fail(c, err, e.CodeInternal)
		return false
	}

	validationError := ginvalidator.Translate(&errs)
	if replace != nil {
		if err := validationError.Replace(replace); err != nil {
			fail(c, err, e.CodeInternal)
			return false
		}
	}
	var fieldViolations []*e.BadRequestFieldViolation
	for field, desc := range validationError {
		fieldViolation := &e.BadRequestFieldViolation{
			Field:       field,
			Description: desc,
		}
		fieldViolations = append(fieldViolations, fieldViolation)
	}
	badRequest := &e.BadRequest{FieldViolations: fieldViolations}
	fail(c, validationError, e.CodeInvalidArgument, badRequest)
	return false
}
