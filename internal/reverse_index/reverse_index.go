package reverseindex

import (
	"search/pb"
)

type IReverseIndex interface {
	Add(doc *pb.Document)               // 将文档加入进倒排索引
	Remove(doc *pb.Document)            // 删除索引中的某文档
	Search(term *pb.TermQuery) []string // 搜索
}
