package config_test

import (
	"github.com/stretchr/testify/assert"
	"project/app/pkg/config"
	"testing"
)

func TestNewViper(t *testing.T) {
	a := assert.New(t)
	v, err := config.NewViper(
		"./testdata/config.yaml",
		"./testdata/secret.yaml",
	)
	a.Nil(err)
	a.Equal("v2", v.GetString("version"))
	a.Equal("/root/a.log", v.GetString("logpath"))
	a.Equal("root", v.GetString("mysql.user"))
	a.Equal("mysql-secret", v.GetString("mysql.password"))
	a.Equal("aliyun-secret-key", v.GetString("aliyun-secret-key"))
}
