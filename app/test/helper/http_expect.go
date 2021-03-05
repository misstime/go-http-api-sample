// 本包定义：用于测试的助手函数

package helper

import (
	"github.com/gavv/httpexpect/v2"
	"net/http"
	"testing"
)

func NewHttpExcept(t *testing.T, handler http.Handler) *httpexpect.Expect {
	t.Helper()
	return httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(handler),
			Jar:       httpexpect.NewJar(),
		},
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})
}
