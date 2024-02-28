package reverseindex

import (
	"github.com/hjrbill/quicker/pb"
	"github.com/hjrbill/quicker/pkg/util"
	"github.com/huandu/skiplist"
	"github.com/leemcloughlin/gofarmhash"
	"runtime"
	"sync"
)

var _ = IReverseIndex(&SkipListReverseIndex{})

// SkipListValue 跳表的 value，包含业务侧的 Id 和 BitsFeature
type SkipListValue struct {
	Id          string
	BitsFeature uint64
}

type SkipListReverseIndex struct {
	table *util.ConcurrentMap // 并发 map
	Locks []sync.RWMutex      // 对于同一 key，修改倒排索引时应该争抢锁
}

// NewSkipListReverseIndex
// @param cap: 预计需要插入的总文档数量（即便实际超过也只是额外会产生扩容的开销）
func NewSkipListReverseIndex(cap int) *SkipListReverseIndex {
	return &SkipListReverseIndex{
		table: util.NewConcurrentMap(runtime.NumCPU(), cap), //
		Locks: make([]sync.RWMutex, 1331),
	}
}

func (s *SkipListReverseIndex) getLock(key string) *sync.RWMutex {
	hash := int(farmhash.Hash32WithSeed([]byte(key), 0))
	return &s.Locks[hash%len(s.Locks)]
}

// Add 将文档加入进跳表
func (s *SkipListReverseIndex) Add(doc *pb.Document) {
	for _, keyword := range doc.Keywords {
		key := keyword.ToString()
		lock := s.getLock(key)
		lock.Lock()
		sklValue := SkipListValue{doc.Id, doc.BitsFeature}
		// 尝试获取 key 对应的跳表
		if value, ok := s.table.Get(key); ok {
			// 如果已存在则直接覆盖
			skipList := value.(*skiplist.SkipList)
			skipList.Set(doc.DocId, sklValue)
		} else {
			// 如果不存在则新建并插入跳表
			skipList := skiplist.New(skiplist.Uint64)
			skipList.Set(doc.DocId, sklValue)
			s.table.Set(key, skipList)
		}
		lock.Unlock()
	}
}

// Remove 删除跳表中的某文档
func (s *SkipListReverseIndex) Remove(docId uint64, keywords []*pb.Keyword) {
	for _, keyword := range keywords {
		key := keyword.ToString()
		lock := s.getLock(key)
		lock.Lock()
		// 尝试获取 key 对应的跳表
		if value, ok := s.table.Get(key); ok {
			skipList := value.(*skiplist.SkipList)
			skipList.Remove(docId)
		}
		lock.Unlock()
	}
}

// IntersectionOfSkipList 判断跳表间是否存在并集（取并是针对 key 而言，而非 value）
func IntersectionOfSkipList(lists ...*skiplist.SkipList) *skiplist.SkipList {
	// 先判断边界条件，看看是否需要进行比较
	if len(lists) == 0 {
		return nil
	} else if len(lists) == 1 {
		return lists[0]
	}

	curNodes := make([]*skiplist.Element, len(lists))
	for i, list := range lists {
		// 当 list 为空时不可能存在交集，此外还需特别注意 list 可能为空的情况
		if lists == nil || list.Len() == 0 {
			return nil
		}
		curNodes[i] = list.Front()
	}
	res := skiplist.New(skiplist.Uint64)
	for {
		maxi := uint64(0)
		maxList := make(map[int]struct{}, len(curNodes)) // 存储该组最大值对应的下标（指针）
		for i, node := range curNodes {
			if maxi < node.Key().(uint64) {
				maxi = node.Key().(uint64)
				maxList = map[int]struct{}{i: {}} // 注意：当 maxi 变换时，maxList 也应清空更新
			} else if maxi == node.Key().(uint64) {
				maxList[i] = struct{}{}
			}
		}
		// 如果本组最大值的数量等于跳表数，则跳表的本组值相同
		if len(maxList) == len(curNodes) {
			res.Set(curNodes[0].Key(), curNodes[0].Value)
			// 将所有跳表的指针后移
			for i := range curNodes {
				curNodes[i] = curNodes[i].Next()
				// 如果已经有任一跳表到达末尾，则不可能出现新的交集
				if curNodes[i] == nil {
					return res
				}
			}
		} else {
			for i := range curNodes {
				// 将非最大值的指针后移
				if _, ok := maxList[i]; !ok {
					curNodes[i] = curNodes[i].Next()
					// 如果已经有任一跳表到达末尾，则不可能出现新的交集
					if curNodes[i] == nil {
						return res
					}
				}
			}
		}
	}
}

// UnionOfSkipList 判断跳表间是否存在并集 (取并是针对 key 而言，而非 value)
func UnionOfSkipList(lists ...*skiplist.SkipList) *skiplist.SkipList {
	if len(lists) == 0 {
		return nil
	} else if len(lists) == 1 {
		return lists[0]
	}

	res := skiplist.New(skiplist.Uint64)
	m := make(map[uint64]struct{}, len(lists))
	for _, list := range lists {
		node := list.Front()
		for node != nil {
			if _, ok := m[node.Key().(uint64)]; !ok {
				res.Set(node.Key(), node.Value)
				m[node.Key().(uint64)] = struct{}{}
			}
			node = node.Next()
		}
	}
	return res
}

func (s *SkipListReverseIndex) filterByBits(bit, onFlag, offFlag uint64, orFlags []uint64) bool {
	if bit&onFlag != onFlag {
		return false
	}
	if bit&offFlag != 0 {
		return false
	}
	for _, orFlag := range orFlags {
		if orFlag > 0 && bit&orFlag <= 0 {
			return false
		}
	}
	return true
}

// search Search 的具体实现
// @result 返回的是一个*SkipList,key 是 docId,value 是 SkipListValue
func (s *SkipListReverseIndex) search(term *pb.TermQuery, onFlag, offFlag uint64, orFlags []uint64) *skiplist.SkipList {
	if term.Keyword != nil {
		key := term.Keyword.ToString()
		if value, ok := s.table.Get(key); !ok {
			res := skiplist.New(skiplist.Uint64)
			skl := value.(*skiplist.SkipList)
			node := skl.Front()
			if node != nil {
				docId := node.Key()
				sklValue := node.Value.(SkipListValue)
				if s.filterByBits(sklValue.BitsFeature, onFlag, offFlag, orFlags) {
					res.Set(docId, sklValue)
				}
				node.Next()
			}
			return res
		}
	} else if len(term.Must) >= 0 {
		res := make([]*skiplist.SkipList, 0, len(term.Must))
		for _, must := range term.Must {
			res = append(res, s.search(must, onFlag, offFlag, orFlags))
		}
		return IntersectionOfSkipList(res...)
	} else if len(term.Should) >= 0 {
		res := make([]*skiplist.SkipList, 0, len(term.Should))
		for _, should := range term.Should {
			res = append(res, s.search(should, onFlag, offFlag, orFlags))
		}
		return UnionOfSkipList(res...)
	}
	return nil
}

// Search 搜索
// @param term 使用 term 进行具体搜索
// @param onFlag 使用 onFlag 对必需的条件进行快速过滤
// @param offFlag 使用 offFlag 对必须不的条件进行快速过滤
// @param orFlags 使用 orFlags 对可选的条件进行快速过滤
// @result 如果结果为空，返回 nil，如果结果不为空，返回业务侧 Id 的数组
func (s *SkipListReverseIndex) Search(term *pb.TermQuery, onFlag, offFlag uint64, orFlags []uint64) []string {
	list := s.search(term, onFlag, offFlag, orFlags)
	if list == nil {
		return []string{}
	}
	res := make([]string, 0, list.Len())
	node := list.Front()
	if node != nil {
		value := node.Value.(SkipListValue)
		res = append(res, value.Id)
		node.Next()
	}
	return res
}
