package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gogo/protobuf/proto"
	model "github.com/hjrbill/quicker/demo/gen/video"
	"github.com/hjrbill/quicker/demo/job"
	"github.com/hjrbill/quicker/demo/job/search"
	"github.com/hjrbill/quicker/demo/job/search/common"
	"github.com/hjrbill/quicker/demo/param"
	pb "github.com/hjrbill/quicker/gen"
	"github.com/hjrbill/quicker/index_service"
	"log"
	"net/http"
	"strings"
)

var Indexer index_service.IIndexer

func clearnKeywords(words []string) []string {
	keywords := make([]string, 0, len(words))
	for _, w := range words {
		word := strings.TrimSpace(strings.ToLower(w))
		if len(word) > 0 {
			keywords = append(keywords, word)
		}
	}
	return keywords
}

// Search 搜索接口
func Search(ctx *gin.Context) {
	var request param.SearchRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		log.Printf("bind request parameter failed: %s", err)
		ctx.String(http.StatusBadRequest, "invalid json")
		return
	}

	keywords := clearnKeywords(request.Keywords)
	if len(keywords) == 0 && len(request.Author) == 0 {
		ctx.String(http.StatusBadRequest, "关键词和作者不能同时为空")
		return
	}
	query := new(pb.TermQuery)
	if len(keywords) > 0 {
		for _, word := range keywords {
			query = query.And(pb.NewTermQuery("content", word)) //满足关键词
		}
	}
	if len(request.Author) > 0 {
		query = query.And(pb.NewTermQuery("author", strings.ToLower(request.Author))) //满足作者
	}
	orFlags := []uint64{job.GetClassBits(request.Classes)} //满足类别
	docs, err := Indexer.Search(query, 0, 0, orFlags)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "search failed")
		return
	}
	videos := make([]model.Video, 0, len(docs))
	for _, doc := range docs {
		var video model.Video
		if err := proto.Unmarshal(doc.Bytes, &video); err == nil {
			if video.View >= int32(request.ViewFrom) && (request.ViewTo <= 0 || video.View <= int32(request.ViewTo)) { //满足播放量的范围
				videos = append(videos, video)
			}
		}
	}
	log.Printf("return %d videos", len(videos))
	ctx.JSON(http.StatusOK, videos) //把搜索结果以 json 形式返回给前端
}

// SearchAll 搜索全站视频
func SearchAll(ctx *gin.Context) {
	var request param.SearchRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		log.Printf("bind request parameter failed: %s", err)
		ctx.String(http.StatusBadRequest, "invalid json")
		return
	}

	request.Keywords = clearnKeywords(request.Keywords)
	if len(request.Keywords) == 0 && len(request.Author) == 0 {
		ctx.String(http.StatusBadRequest, "关键词和作者不能同时为空")
		return
	}

	searchCtx := &common.VideoSearchContext{
		Ctx:     context.Background(),
		Request: &request,
		Indexer: Indexer,
	}
	searcher := search.NewAllVideoSearcher()
	videos := searcher.Search(searchCtx)

	ctx.JSON(http.StatusOK, videos) //把搜索结果以 json 形式返回给前端
}

// SearchByAuthor up 主在后台搜索自己的视频
func SearchByAuthor(ctx *gin.Context) {
	var request param.SearchRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		log.Printf("bind request parameter failed: %s", err)
		ctx.String(http.StatusBadRequest, "invalid json")
		return
	}

	request.Keywords = clearnKeywords(request.Keywords)
	if len(request.Keywords) == 0 {
		ctx.String(http.StatusBadRequest, "关键词不能为空")
		return
	}

	userName, ok := ctx.Value("user_name").(string) //从 gin.Context 里取得 userName
	if !ok || len(userName) == 0 {
		ctx.String(http.StatusBadRequest, "获取不到登录用户名")
		return
	}
	searchCtx := &common.VideoSearchContext{
		Ctx:     context.WithValue(context.Background(), common.UN("user_name"), userName), //把 userName 放到 context 里
		Request: &request,
		Indexer: Indexer,
	}
	searcher := search.NewAuthorVideoSearcher()
	videos := searcher.Search(searchCtx)

	ctx.JSON(http.StatusOK, videos) //把搜索结果以 json 形式返回给前端
}
