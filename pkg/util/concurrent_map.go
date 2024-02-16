package util

import (
	"github.com/leemcloughlin/gofarmhash"
	"sync"
)

type MapEntry struct {
	Key   string
	value any
}

// MapIterator 引入迭代器模式，实现对自定义 map 的迭代遍历
type MapIterator interface {
	Next() *MapEntry
}

var _ MapIterator = (*ConcurrentMapIterator)(nil) // 对 ConcurrentMapIterator 是否实现了对 MapIterator 的继承进行检查

type ConcurrentMap struct {
	maps  []map[string]any
	seg   int            // 维护 maps 的大小
	locks []sync.RWMutex // 对于每一个 map 分配一个专用的锁，避免共用锁导致的占用等待问题
	seed  uint32
}

// NewConcurrentMap
// @param seg 内部包含几个 map
// @param cap 预估 map 中一共容纳多少元素（即便实际超过也只是会额外产生扩容的开销）
func NewConcurrentMap(seg, cap int) *ConcurrentMap {
	maps := make([]map[string]any, seg)
	locks := make([]sync.RWMutex, seg)
	for i := range maps {
		maps[i] = make(map[string]any, cap/seg)
	}
	return &ConcurrentMap{
		maps:  maps,
		seg:   seg,
		locks: locks,
		seed:  0,
	}
}

// 判断 key 对应到哪个 map
func (m *ConcurrentMap) getSegIndex(key string) int {
	hash := int(farmhash.Hash32WithSeed([]byte(key), m.seed))
	return hash % m.seg // 对槽位取模
}

func (m *ConcurrentMap) Get(key string) (any, bool) {
	index := m.getSegIndex(key)
	m.locks[index].RLock()
	defer m.locks[index].RUnlock()
	if v, ok := m.maps[index][key]; ok {
		return v, true
	}
	return nil, false
}

func (m *ConcurrentMap) Set(key string, value any) {
	index := m.getSegIndex(key)
	m.locks[index].Lock()
	defer m.locks[index].Unlock()
	m.maps[index][key] = value
}

// ConcurrentMapIterator 实现 ConcurrentMap 的迭代器
type ConcurrentMapIterator struct {
	m   *ConcurrentMap
	key [][]string // row 为 maps 数组的下标，col 为该下标对应的 map 中 key 的下标
	row int        // 维护 key 的 row
	col int        // 维护 key 的 col
}

func (i ConcurrentMapIterator) Next() *MapEntry {
	// 如果当前已经完成遍历，直接返回 nil
	if i.row > len(i.key)-1 {
		return nil
	}
	// 如果本行为空，递归的跳行，直到找到非空行
	if len(i.key[i.row]) == 0 {
		i.row++
		return i.Next()
	}
	// 如果已经遍历完本列，维护下标并递归
	if i.col > len(i.key[i.row])-1 {
		i.row++
		i.col = 0
		return i.Next()
	}

	key := i.key[i.row][i.col]
	value, ok := i.m.Get(key)
	// 先完成对下标的维护
	if i.col >= len(i.key[i.row])-1 {
		i.row++
		i.col = 0
	} else {
		i.col++
	}
	// 如果没有取到值，递归到下一个
	if !ok {
		return i.Next()
	}
	// 否则返回取到的值
	return &MapEntry{
		Key:   key,
		value: value,
	}
}
