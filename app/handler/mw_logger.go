package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"project/app/pkg/config"
	"time"
)

// LoggerMiddleware 用于记录 http 日志，包含请求信息、响应信息、以及错误。
type LoggerMiddleware struct {
	zapLogger *zap.Logger
}

//func NewLoggerMiddleware(isDebug config.IsDebug) (mw *LoggerMiddleware, cleanup func(), err error) {
//	var zapLogger *zap.Logger
//
//	if isDebug {
//		zapLogger, err = zap.NewDevelopment()
//	} else {
//		zapLogger, err = zap.NewProduction()
//	}
//
//	if err != nil {
//		cleanup = func() {
//			if err := zapLogger.Sync(); err != nil {
//				fmt.Printf("http logger middleware zap sync() error: %s", err)
//			}
//		}
//	}
//
//	return &LoggerMiddleware{zapLogger: zapLogger}, cleanup, err
//}

func NewLoggerMiddleware(zapLogger *zap.Logger) *LoggerMiddleware {
	return &LoggerMiddleware{zapLogger: zapLogger}
}

// CreateGinHandler 生成 gin.HandlerFunc 实例 - 即：gin middleware 函数
func (mw *LoggerMiddleware) CreateGinHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		// 获取 gin.Context 中附加的三个数据：
		// - response body
		// - api error：api 请求过程中发生的错误
		// - log level：zap 记录日志时使用的错误级别
		iBody, exists := c.Get(contextKeyBody)
		if !exists {
			fmt.Printf("http logger context key `%s` not exists\n", contextKeyBody)
			return
		}
		body, ok := iBody.(*body)
		if !ok {
			fmt.Printf("http logger context key `%s` type invalid\n", contextKeyBody)
			return
		}
		iLogLevel, exists := c.Get(contextKeyLogLevel)
		if !exists {
			fmt.Printf("http logger context key `%s` not exists\n", contextKeyLogLevel)
			return
		}
		logLevel, ok := iLogLevel.(zapcore.Level)
		if !ok {
			fmt.Printf("http logger context key `%s` type invalid\n", contextKeyLogLevel)
			return
		}
		apiError, _ := c.Get(contextKeyError)

		endtime := time.Now()

		var zapFields []zap.Field
		zapFields = append(
			zapFields,
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("http_status", c.Writer.Status()),
			zap.Uint8("code", uint8(body.Code)),
			zap.String("code_name", body.Status),
			zap.String("msg", body.Message),
			zap.Any("error", apiError),
			zap.Time("start_time", startTime),
			zap.Time("end_time", endtime),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.String("referer", c.Request.Referer()),
		)
		switch logLevel {
		case zapcore.DebugLevel:
			mw.zapLogger.Debug(body.Status, zapFields...)
		case zapcore.InfoLevel:
			mw.zapLogger.Debug(body.Status, zapFields...)
		case zapcore.WarnLevel:
			mw.zapLogger.Debug(body.Status, zapFields...)
		case zapcore.FatalLevel:
			mw.zapLogger.Debug(body.Status, zapFields...)
		default:
			fmt.Printf("unsupported log level: %s\n", logLevel)
		}
	}
}

// 实例化一个 *zap.Logger，该实例用于注入 LoggerMiddleware
func NewZapLogger(isDebug config.IsDebug, v *viper.Viper) (zapLogger *zap.Logger, cleanup func(), err error) {
	// @todo 以下为 mock 代码
	if isDebug {
		if zapLogger, err = zap.NewDevelopment(); err != nil {
			return nil, nil, errors.Wrap(err, "new development zapLogger failed")
		}
	} else {
		if zapLogger, err = zap.NewDevelopment(); err != nil {
			return nil, nil, errors.Wrap(err, "new production zapLogger failed")
		}
	}
	cleanup = func() {
		if err := zapLogger.Sync(); err != nil {
			fmt.Printf("zap sync() error: %s", err)
		}
	}
	return zapLogger, cleanup, nil
}
