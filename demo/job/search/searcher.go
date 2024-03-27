package search

import (
	model "github.com/hjrbill/quicker/demo/gen/video"
	"github.com/hjrbill/quicker/demo/job/search/common"
	"github.com/hjrbill/quicker/demo/job/search/filter"
	"github.com/hjrbill/quicker/demo/job/search/recaller"
	"time"

	"log"
	"reflect"
	"sync"
)

type Recaller interface {
	Recall(*common.VideoSearchContext) ([]*model.Video, error)
}

type Filter interface {
	Apply(*common.VideoSearchContext)
}

// VideoSearcher 一个模板方法的超类，规范了 Search 的过程，其子类应该重写 recaller 和 filter
type VideoSearcher struct {
	Recallers []Recaller //实际中，除了正常的关键词召回外，可能还要召回广告
	Filters   []Filter
}

func (searcher *VideoSearcher) WithRecaller(recaller ...Recaller) {
	searcher.Recallers = append(searcher.Recallers, recaller...)
}

func (searcher *VideoSearcher) WithFilter(filter ...Filter) {
	searcher.Filters = append(searcher.Filters, filter...)
}

func (searcher *VideoSearcher) Recall(ctx *common.VideoSearchContext) {
	if len(searcher.Recallers) == 0 {
		return
	}
	wg := sync.WaitGroup{}
	wg.Add(len(searcher.Recallers))
	collection := make(chan *model.Video, 1000)
	// 并发的执行多种召回
	for _, searchRecaller := range searcher.Recallers {
		go func(searchRecaller Recaller) {
			defer wg.Done()
			rule := reflect.TypeOf(searchRecaller).Name()
			videos, err := searchRecaller.Recall(ctx)
			if err != nil {
				log.Printf("召回规则 %s 失败: %s", rule, err)
				return
			}
			log.Printf("以 %s 规则召回了 %d 条视频\n", rule, len(videos))
			for _, video := range videos {
				collection <- video
			}
		}(searchRecaller)
	}
	wg.Wait()
	close(collection)
}

func (searcher *VideoSearcher) Filter(ctx *common.VideoSearchContext) {
	for _, searchFilter := range searcher.Filters { // 顺序执行各过滤规则
		searchFilter.Apply(ctx)
	}
}

func (searcher *VideoSearcher) Search(ctx *common.VideoSearchContext) []*model.Video {
	t1 := time.Now()
	//召回
	searcher.Recall(ctx)
	t2 := time.Now()
	log.Printf("recall %d docs in %d ms", len(ctx.Result), t2.Sub(t1).Milliseconds())
	//过滤
	searcher.Filter(ctx)
	t3 := time.Now()
	log.Printf("after filter remain %d docs in %d ms", len(ctx.Result), t3.Sub(t2).Milliseconds())
	return ctx.Result
}

// 子类
type AllVideoSearcher struct {
	VideoSearcher
}

func NewAllVideoSearcher() *AllVideoSearcher {
	searcher := new(AllVideoSearcher)
	searcher.WithRecaller(&recaller.KeywordRecaller{})
	searcher.WithFilter(&filter.ViewerFilter{})
	return searcher
}

// 子类
type AuthorVideoSearcher struct {
	VideoSearcher
}

func NewAuthorVideoSearcher() *AuthorVideoSearcher {
	searcher := new(AuthorVideoSearcher)
	searcher.WithRecaller(&recaller.KeywordAuthorRecaller{})
	searcher.WithFilter(&filter.ViewerFilter{})
	return searcher
}
