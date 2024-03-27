package index_service

import pb "github.com/hjrbill/quicker/gen"

// IIndexer Sentinel（分布式 grpc 的哨兵）和 Indexer（单机索引）都实现了该接口
type IIndexer interface {
	AddDoc(doc pb.Document) (int, error)
	DeleteDoc(id string) (int, error)
	Search(query *pb.TermQuery, onFlag uint64, offFlag uint64, orFlags []uint64) ([]*pb.Document, error)
	Count() int32
	Close() error
}
