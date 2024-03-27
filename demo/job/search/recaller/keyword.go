package recaller

import (
	"github.com/gogo/protobuf/proto"
	model "github.com/hjrbill/quicker/demo/gen/video"
	"github.com/hjrbill/quicker/demo/job"
	"github.com/hjrbill/quicker/demo/job/search/common"
	"github.com/hjrbill/quicker/pb"
)

type KeywordRecaller struct {
}

func (*KeywordRecaller) Recall(ctx *common.VideoSearchContext) ([]*model.Video, error) {
	request := ctx.Request
	if request == nil {
		return nil, nil
	}
	indexer := ctx.Indexer
	if indexer == nil {
		return nil, nil
	}

	query := new(pb.TermQuery)
	if len(request.Keywords) > 0 {
		for _, keyword := range request.Keywords {
			query.And(pb.NewTermQuery("content", keyword)) // 按关键字过滤
		}
	}
	if len(request.Author) > 0 {
		query.And(pb.NewTermQuery("author", request.Author)) // 按作者过滤
	}
	orFlags := []uint64{job.GetClassBits(request.Classes)} // 按类别过滤
	docs, err := indexer.Search(query, 0, 0, orFlags)
	if err != nil {
		return nil, err
	}

	res := make([]*model.Video, 0, len(docs))
	for _, doc := range docs {
		var video model.Video
		if err := proto.Unmarshal(doc.Bytes, &video); err != nil {
			res = append(res, &video)
		}
	}
	return res, nil
}
