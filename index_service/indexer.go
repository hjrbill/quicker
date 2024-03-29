package index_service

import (
	"bytes"
	"encoding/gob"
	pb "github.com/hjrbill/quicker/gen"
	"github.com/hjrbill/quicker/internal/kvdb"
	reverseindex "github.com/hjrbill/quicker/internal/reverse_index"
	qlog "github.com/hjrbill/quicker/pkg/log"
	"strings"
	"sync/atomic"
)

var _ IIndexer = (*Indexer)(nil)

// Indexer 一个外观，封装正排索引与倒排索引两个模块（子系统）
type Indexer struct {
	forwardIndex kvdb.IKeyValueDB
	reverseIndex reverseindex.IReverseIndex
	docIdCnt     uint64 // 用于记录当前 DocId 的最大值
}

func (indexer *Indexer) Init(cap int, dbType kvdb.DBType, dbPath string) error {
	db, err := kvdb.NewKVDB(dbType, dbPath)
	if err != nil {
		return err
	}
	indexer.forwardIndex = db
	indexer.reverseIndex = reverseindex.NewSkipListReverseIndex(cap)
	return nil
}

func (indexer *Indexer) Close() error {
	err := indexer.forwardIndex.Close()
	if err != nil {
		return err
	}
	return nil
}

// LoadFromForwardIndexFile 从正排索引的数据库文件中加载数据（用于程序重启后）
func (indexer *Indexer) LoadFromForwardIndexFile() int64 {
	reader := bytes.NewReader([]byte{})
	n := indexer.forwardIndex.IterateDB(func(k, v []byte) error {
		reader.Reset(v)
		decoder := gob.NewDecoder(reader)
		var doc pb.Document
		err := decoder.Decode(&doc)
		if err != nil {
			qlog.Warnf("gob decode document failed：%s", err)
			return nil
		}
		indexer.reverseIndex.Add(doc)
		return err
	})
	qlog.Infof("load %d data from forward index %s", n, indexer.forwardIndex.GetDBPath())
	return n
}

func (indexer *Indexer) Count() int32 {
	n := int32(0)
	indexer.forwardIndex.IterateKey(func(key []byte) error {
		n++
		return nil
	})
	return n
}

// AddDoc 在索引上添加文档，如果已存在，则会覆盖
func (indexer *Indexer) AddDoc(doc pb.Document) (int, error) {
	id := strings.TrimSpace(doc.Id)
	if len(id) == 0 {
		return 0, nil
	}
	// 先删除原有的文档
	_, err := indexer.DeleteDoc(id)
	if err != nil {
		return 0, err
	}

	// 生成新的 DocId
	doc.DocId = atomic.AddUint64(&indexer.docIdCnt, 1)
	// 再添加新的文档
	var value bytes.Buffer
	decoder := gob.NewEncoder(&value)
	if err := decoder.Encode(doc); err == nil {
		err := indexer.forwardIndex.Set([]byte(id), value.Bytes())
		if err != nil {
			return 0, err
		}
	} else {
		return 0, err
	}

	indexer.reverseIndex.Add(doc)
	return 1, nil
}

// DeleteDoc 从索引上删除文档
func (indexer *Indexer) DeleteDoc(id string) (int, error) {
	n := 0
	forwardKey := []byte(id)
	//先读正排索引，得到 IntId 和 Keywords
	docBs, err := indexer.forwardIndex.Get(forwardKey)
	if err == nil {
		reader := bytes.NewReader([]byte{})
		if len(docBs) > 0 {
			n = 1
			reader.Reset(docBs)
			decoder := gob.NewDecoder(reader)
			var doc pb.Document
			err := decoder.Decode(&doc)
			if err == nil {
				indexer.reverseIndex.Remove(doc.DocId, doc.Keywords)
			}
		}
	}
	//从正排上删除
	err = indexer.forwardIndex.Delete(forwardKey)
	if err != nil {
		return 0, err
	}
	return n, err
}

func (indexer *Indexer) Search(query *pb.TermQuery, onFlag, offFlag uint64, orFlags []uint64) ([]*pb.Document, error) {
	ids := indexer.reverseIndex.Search(query, onFlag, offFlag, orFlags)
	if len(ids) == 0 {
		return nil, nil
	}

	keys := make([][]byte, 0, len(ids)) // 正排索引的 key
	for _, id := range ids {
		keys = append(keys, []byte(id))
	}
	// 从正排索引中获取序列化后的文档
	docBytes, err := indexer.forwardIndex.BatchGet(keys)
	if err != nil {
		qlog.Warnf("read kv db failed, error: %s", err)
		return nil, err
	}

	result := make([]*pb.Document, 0, len(docBytes))
	reader := bytes.NewReader([]byte{})
	for _, docByte := range docBytes {
		if len(docByte) > 0 {
			reader.Reset(docByte)
			decoder := gob.NewDecoder(reader)
			var doc pb.Document
			err := decoder.Decode(&doc)
			if err == nil { // 如果解码成功，才将文档添加到返回值中
				result = append(result, &doc)
			}
		}
	}
	return result, nil
}
