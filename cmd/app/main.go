package main

import (
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"project/app/handler/pkg/ginvalidator"
	"project/app/pkg/config"
)

var (
	// 模板配置文件路径
	configTemplateFile = flag.String(
		"config_template",
		//"template.yaml",
		"c:/projects/go-http-api-sample/app/config/template.yaml",
		"set config template file which viper will loading.",
	)
	// 机密配置文件路径
	configSecretFile = flag.String(
		"config_secret",
		//"secret.yaml",
		"c:/projects/go-http-api-sample/app/config/secret.yaml",
		"set config secret file which viper will loading.",
	)
)

func main() {
	flag.Parse()

	// 初始化 gin 内部使用的 validator，如：注册自定义验证器、注册翻译器等...
	// 由于要使用到 gin 的 c.ShouldBind() 系列函数，暂时无法通过依赖注入来避免该初始化操作
	if err := ginvalidator.Init(); err != nil {
		panic(errors.Wrap(err, "gin validator initialize failed"))
	}

	app, cleanup, err := CreateApp(config.FilePath(*configTemplateFile), config.FilePath(*configSecretFile))
	if err != nil {
		panic(err)
	} else {
		defer cleanup()
		fmt.Println(app.Run())
	}
	// ./app.exe -config_template=C:/projects/go-http-api-sample/app/config/template.yaml --config_secret=C:/projects/go-http-api-sample/app/config/secret.yaml
}
