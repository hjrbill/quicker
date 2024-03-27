package search

import (
	model "github.com/hjrbill/quicker/demo/gen/video"
	"github.com/hjrbill/quicker/demo/job/search/common"
	"github.com/hjrbill/quicker/demo/job/search/filter"
	"github.com/hjrbill/quicker/demo/job/search/recaller"
	"golang.org/x/exp/maps"
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
	//并行执行多路召回
	collection := make(chan *model.Video, 1000)
	wg := sync.WaitGroup{}
	wg.Add(len(searcher.Recallers))
	for _, recaller := range searcher.Recallers {
		go func(recaller Recaller) {
			defer wg.Done()
			rule := reflect.TypeOf(recaller).Name()
			result, err := recaller.Recall(ctx)
			if err != nil {
				log.Printf("recall %s failed: %v", rule, err)
				return
			}
			log.Printf("recall %d talents by %s", len(result), rule)
			for _, video := range result {
				collection <- video
			}
		}(recaller)
	}
	//通过 map 合并多路召回的结果
	videoMap := make(map[string]*model.Video, 1000)
	receiveFinish := make(chan struct{})
	go func() {
		for {
			video, ok := <-collection
			if !ok {
				break
			}
			videoMap[video.ID] = video
		}
		receiveFinish <- struct{}{}
	}()
	wg.Wait()
	close(collection)
	<-receiveFinish

	ctx.Result = maps.Values(videoMap)
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
	searcher.WithRecaller(recaller.KeywordRecaller{})
	searcher.WithFilter(filter.ViewerFilter{})
	return searcher
}

// 子类
type AuthorVideoSearcher struct {
	VideoSearcher
}

func NewAuthorVideoSearcher() *AuthorVideoSearcher {
	searcher := new(AuthorVideoSearcher)
	searcher.WithRecaller(recaller.KeywordAuthorRecaller{})
	searcher.WithFilter(filter.ViewerFilter{})
	return searcher
}
