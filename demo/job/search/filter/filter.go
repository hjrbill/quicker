package filter

import (
	model "github.com/hjrbill/quicker/demo/gen/video"
	"github.com/hjrbill/quicker/demo/job/search/common"
)

type Filter interface {
	Apply(*common.VideoSearchContext)
}

var _ Filter = (*ViewerFilter)(nil)

type ViewerFilter struct {
}

// Apply 根据播放量过滤视频（这是业务方的特殊需求，应该由业务方实现）
func (ViewerFilter) Apply(ctx *common.VideoSearchContext) {
	request := ctx.Request
	if request == nil {
		return
	}
	if request.ViewFrom >= request.ViewTo {
		return
	}
	vidoes := make([]*model.Video, 0, len(ctx.Result))
	for _, video := range ctx.Result {
		if video.View >= int32(request.ViewFrom) && video.View <= int32(request.ViewTo) {
			vidoes = append(vidoes, video)
		}
	}
	ctx.Result = vidoes
}
