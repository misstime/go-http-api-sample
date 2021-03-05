// +build wireinject

package main

import (
	"github.com/google/wire"
	"project/app"
	"project/app/handler"
	"project/app/pkg/cache"
	"project/app/pkg/config"
	"project/app/pkg/sms"
	"project/app/service"
)

var providerSet = wire.NewSet(
	// 公共 provider
	config.NewViper,
	config.NewIsDebug,
	cache.NewGoCache,

	// app
	app.NewApp,
	app.NewHttpAddresses,

	// LoggerMiddleware
	handler.NewLoggerMiddleware,
	handler.NewZapLogger,

	// RecoveryMiddleware
	wire.Value(&handler.RecoveryMiddleware{}),

	// LoginSmsCtrl
	handler.NewLoginSmsCtrl,
	service.NewLoginSmsService,
	wire.Bind(new(service.ISms), new(*service.LoginSmsService)),
	sms.NewAliyunLoginSms,
	wire.Bind(new(sms.Sender), new(*sms.AliyunLoginSms)),
)

func CreateApp(configFiles ...config.FilePath) (*app.App, func(), error) {
	panic(wire.Build(providerSet))
}
