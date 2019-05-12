package media

import "github.com/JefferyWang/workwechat/context"

// Media 素材管理
type Media struct {
	*context.Context
}

// NewMedia 实例化
func NewMedia(context *context.Context) *Media {
	media := new(Media)
	media.Context = context
	return media
}
