package cache

import (
	"github.com/patrickmn/go-cache"
	"time"
)

const (
	// go-cache 缓存默认过期时间
	goCacheDefaultExpiration = 30 * time.Minute
	// go-cache 缓存清理时间
	goCacheCleanupInterval = 60 * time.Minute
)

// NewGoCache 实例化一个 go-cache 库的缓存客户端
func NewGoCache() *cache.Cache {
	return cache.New(goCacheDefaultExpiration, goCacheCleanupInterval)
}
