package context

import (
	"net/http"
	"sync"

	"github.com/JefferyWang/workwechat/cache"
)

// Context struct
type Context struct {
	CorpID     string // 企业ID
	CorpSecret string // 应用的凭证密钥

	Cache cache.Cache

	Writer  http.ResponseWriter
	Request *http.Request

	//accessTokenLock 读写锁 同一个AppID一个
	accessTokenLock *sync.RWMutex
}

// Query returns the keyed url query value if it exists
func (ctx *Context) Query(key string) string {
	value, _ := ctx.GetQuery(key)
	return value
}

// GetQuery is like Query(), it returns the keyed url query value
func (ctx *Context) GetQuery(key string) (string, bool) {
	req := ctx.Request
	if values, ok := req.URL.Query()[key]; ok && len(values) > 0 {
		return values[0], true
	}
	return "", false
}
