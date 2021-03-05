package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type FilePath string

// NewViper 实例化一个 *viper.Viper
// 当传入多个文件路径时，获取的配置信息会依次向前覆盖，越靠后优先级越高。
func NewViper(paths ...FilePath) (*viper.Viper, error) {
	if len(paths) == 0 {
		return nil, errors.New("at least one configuration is required")
	}

	v := viper.New()

	for k, path := range paths {
		v.SetConfigFile(string(path))
		if k == 0 {
			if err := v.ReadInConfig(); err != nil {
				return nil, err
			}
		} else {
			if err := v.MergeInConfig(); err != nil {
				return nil, err
			}
		}
	}

	return v, nil
}

type IsDebug bool

// NewIsDebug 判断是否处于开发者模式
func NewIsDebug(v *viper.Viper) IsDebug {
	return IsDebug(v.GetBool("isDebug"))
}
