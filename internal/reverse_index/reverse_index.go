package reverseindex

import (
	"search/pb"
)

type IReverseIndex interface {
	Add(doc *pb.Document)                                                         // 将文档加入进倒排索引
	Remove(docId uint64, keywords []*pb.Keyword)                                  // 删除索引中的某文档
	Search(term *pb.TermQuery, onFlag, offFlag uint64, orFlags []uint64) []string // 搜索，如果结果为空，返回 nil，如果结果不为空，返回业务方 Id 数组
}
