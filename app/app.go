package app

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"project/app/handler"
	"project/app/handler/pkg/ginvalidator"
	"project/app/pkg/config"
)

type App struct {
	isDebug       config.IsDebug
	httpAddresses []string // http 监听地址

	loggerMiddleware   *handler.LoggerMiddleware   // http 日志中间件
	recoveryMiddleware *handler.RecoveryMiddleware // recovery 中间件

	loginSmsCtrl *handler.LoginSmsCtrl // 登录验证码控制器
}

type HttpAddresses []string

func NewHttpAddresses(v *viper.Viper) HttpAddresses {
	return v.GetStringSlice("addr")
}

func NewApp(
	isDebug config.IsDebug,
	httpAddresses HttpAddresses,

	loggerMiddleware *handler.LoggerMiddleware,
	recoveryMiddleware *handler.RecoveryMiddleware,

	loginSmsCtrl *handler.LoginSmsCtrl,
) *App {
	return &App{
		isDebug:            isDebug,
		httpAddresses:      httpAddresses,
		loggerMiddleware:   loggerMiddleware,
		recoveryMiddleware: recoveryMiddleware,
		loginSmsCtrl:       loginSmsCtrl,
	}
}

func (app *App) Run() error {
	if !app.isDebug {
		gin.SetMode(gin.ReleaseMode)
	}
	if err := ginvalidator.Init(); err != nil {
		return err
	}

	engine := gin.New()
	r := engine.Use(
		app.loggerMiddleware.CreateGinHandler(),
		app.recoveryMiddleware.CreateGinHandler(),
	)

	// 手机验证码发送模块（聚合所有的手机验证码发送操作）
	r = engine.Group("/sms")
	{
		// 发送登录手机验证码
		r.POST("/login", app.loginSmsCtrl.Send)
	}

	return engine.Run(app.httpAddresses...)
}
