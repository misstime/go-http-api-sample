package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net"
	"net/http"
	"os"
	"project/app/handler/pkg/e"
	"strings"
)

type RecoveryMiddleware struct {
}

func (mw *RecoveryMiddleware) CreateGinHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				// If the connection is dead, we can't write a status to it.
				if brokenPipe {
					switch err.(type) {
					case string:
						fail(c, errors.New("application panic: "+err.(string)), e.CodeInternal)
					case error:
						fail(c, errors.Wrap(err.(error), "application panic"), e.CodeInternal)
					}
				} else {
					c.AbortWithStatus(http.StatusInternalServerError)
				}
			}
		}()
		c.Next()
	}
}
