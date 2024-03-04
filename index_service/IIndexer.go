package index_service

import "github.com/hjrbill/quicker/pb"

type IIndexer interface {
	AddDoc(doc pb.Document) (int, error)
	DeleteDoc(id string) (int, error)
	Search(query *pb.TermQuery, onFlag uint64, offFlag uint64, orFlags []uint64) ([]*pb.Document, error)
	Count() int32
	Close() error
}
