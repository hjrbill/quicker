package common

import (
	"context"
	model "github.com/hjrbill/quicker/demo/gen/video"
	"github.com/hjrbill/quicker/demo/param"
	"github.com/hjrbill/quicker/index_service"
)

type VideoSearchContext struct {
	Ctx     context.Context
	Indexer index_service.IIndexer
	Request *param.SearchRequest
	Result  []*model.Video
}

type UN string
