package recaller

import (
	"errors"
	"github.com/gogo/protobuf/proto"
	model "github.com/hjrbill/quicker/demo/gen/video"
	"github.com/hjrbill/quicker/demo/job"
	"github.com/hjrbill/quicker/demo/job/search/common"
	"github.com/hjrbill/quicker/pb"
	"strings"
)

type KeywordRecaller struct {
}

func (KeywordRecaller) Recall(ctx *common.VideoSearchContext) ([]*model.Video, error) {
	request := ctx.Request
	if request == nil {
		return nil, errors.New("request is nil")
	}
	indexer := ctx.Indexer
	if indexer == nil {
		return nil, errors.New("indexer is nil")
	}
	keywords := request.Keywords
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
	docs, err := indexer.Search(query, 0, 0, orFlags)
	if err != nil {
		return nil, err
	}
	videos := make([]*model.Video, 0, len(docs))
	for _, doc := range docs {
		var video model.Video
		if err := proto.Unmarshal(doc.Bytes, &video); err == nil {
			videos = append(videos, &video)
		}
	}
	return videos, nil
}
