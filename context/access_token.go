package context

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/JefferyWang/workwechat/util"
)

const (
	// accessTokenURL 获取access_token的接口
	accessTokenURL = "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s"
)

// AccessTokenResp struct
type AccessTokenResp struct {
	util.CommonError

	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

func (ctx *Context) getAccessTokenCacheKey() string {
	return fmt.Sprintf("workwechat_access_token_%s", ctx.CorpSecret)
}

// SetAccessTokenLock 设置读写锁（一个appID一个读写锁）
func (ctx *Context) SetAccessTokenLock(l *sync.RWMutex) {
	ctx.accessTokenLock = l
}

// GetAccessToken 获取access_token
func (ctx *Context) GetAccessToken() (accessToken string, err error) {
	ctx.accessTokenLock.Lock()
	defer ctx.accessTokenLock.Unlock()

	val := ctx.Cache.Get(ctx.getAccessTokenCacheKey())
	if val != nil {
		accessToken = val.(string)
		return
	}

	//从微信服务器获取
	var accessTokenResp AccessTokenResp
	accessTokenResp, err = ctx.GetAccessTokenFromServer()
	if err != nil {
		return
	}

	accessToken = accessTokenResp.AccessToken
	return
}

// GetAccessTokenFromServer 强制从企业微信服务器获取token
func (ctx *Context) GetAccessTokenFromServer() (accessTokenResp AccessTokenResp, err error) {
	log.Printf("GetAccessTokenFromServer")
	url := fmt.Sprintf(accessTokenURL, ctx.CorpID, ctx.CorpSecret)
	var body []byte
	body, err = util.HTTPGet(url)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &accessTokenResp)
	if err != nil {
		return
	}
	if accessTokenResp.ErrCode != 0 {
		err = fmt.Errorf("get work wechat access_token error : errcode=%v , errormsg=%v", accessTokenResp.ErrCode, accessTokenResp.ErrMsg)
		return
	}

	expires := accessTokenResp.ExpiresIn - 1500
	err = ctx.Cache.Set(ctx.getAccessTokenCacheKey(), accessTokenResp.AccessToken, time.Duration(expires)*time.Second)
	return
}
