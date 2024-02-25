package kvdb

import (
	"errors"
	"log"
	"os"
	"path"
	"sync/atomic"

	"github.com/dgraph-io/badger/v4"
)

var _ = IKeyValueDB(&BadgerDB{})

type BadgerDB struct {
	db   *badger.DB
	path string
}

func (s *BadgerDB) WithPath(path string) *BadgerDB {
	s.path = path
	return s
}

func (s *BadgerDB) Open() error {
	DataDir := s.GetDBPath()
	if err := os.MkdirAll(path.Dir(DataDir), os.ModePerm); err != nil { //如果 DataDir 对应的文件夹已存在则什么都不做，如果 DataDir 对应的文件已存在则返回错误
		return err
	}
	option := badger.DefaultOptions(DataDir).WithNumVersionsToKeep(1).WithLoggingLevel(badger.ERROR)
	db, err := badger.Open(option) //文件只能被一个进程使用，如果不调用 Close 则下次无法 Open（手动释放锁的办法：把 LOCK 文件删掉）
	if err != nil {
		return err
	} else {
		s.db = db
		return nil
	}
}

// Close 把内存中的数据 flush 到磁盘，同时释放文件锁 (如果没有 close，再 open 时会丢失很多数据)
func (s *BadgerDB) Close() error {
	return s.db.Close()
}

func (s *BadgerDB) GetDBPath() string {
	return s.path
}

// CheckAndGC 执行 GC
func (s *BadgerDB) CheckAndGC() {
	lsmSize1, vlogSize1 := s.db.Size()
	for {
		if err := s.db.RunValueLogGC(0.5); errors.Is(err, badger.ErrNoRewrite) || errors.Is(err, badger.ErrRejected) {
			break
		}
	}
	lsmSize2, vlogSize2 := s.db.Size()
	if vlogSize2 < vlogSize1 {
		// TODO：log 应该被替换
		log.Printf("badger before GC, LSM %d, vlog %d. after GC, LSM %d, vlog %d", lsmSize1, vlogSize1, lsmSize2, vlogSize2)
	} else {
		log.Printf("collect zero garbage")
	}
}

// Set 为单个写操作开一个事务
func (s *BadgerDB) Set(k, v []byte) error {
	err := s.db.Update(func(txn *badger.Txn) error { //db.Update 相当于打开了一个读写事务:db.NewTransaction(true)。用 db.Update 的好处在于不用显式调用 Txn.Commit() 了
		//duration := time.Hour * 87600
		return txn.Set(k, v) //duration 是能存活的时长
	})
	return err
}

// BatchSet 多个写操作使用一个事务
func (s *BadgerDB) BatchSet(keys, values [][]byte) error {
	if len(keys) != len(values) {
		return errors.New("key value not the same length")
	}
	var err error
	txn := s.db.NewTransaction(true)
	for i, key := range keys {
		value := values[i]
		if err = txn.Set(key, value); err != nil {
			_ = txn.Commit() //发生异常时就提交老事务，然后开一个新事务，重试 set
			txn = s.db.NewTransaction(true)
			_ = txn.Set(key, value)
		}
	}
	txn.Commit()
	return err
}

// Get 如果 key 不存在会返回 error:Key not found
func (s *BadgerDB) Get(k []byte) ([]byte, error) {
	var ival []byte
	err := s.db.View(func(txn *badger.Txn) error { //db.View 相当于打开了一个读写事务:db.NewTransaction(true)。用 db.Update 的好处在于不用显式调用 Txn.Discard() 了
		item, err := txn.Get(k)
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			ival = val
			return nil
		})
		return err
	})
	if len(ival) == 0 {
		return nil, ErrDataNotFound
	}
	return ival, err
}

// BatchGet 返回的 values 与传入的 keys 顺序保持一致。如果 key 不存在或读取失败则对应的 value 是空数组
func (s *BadgerDB) BatchGet(keys [][]byte) ([][]byte, error) {
	var err error
	txn := s.db.NewTransaction(false) //只读事务
	values := make([][]byte, len(keys))
	for i, key := range keys {
		var item *badger.Item
		item, err = txn.Get(key)
		if err == nil {
			var ival []byte
			err = item.Value(func(val []byte) error {
				ival = val
				return nil
			})
			if err == nil {
				values[i] = ival
			} else { //拷贝失败
				values[i] = []byte{} //拷贝失败就把 value 设为空数组
			}
		} else { //读取失败
			values[i] = []byte{}                        //读取失败就把 value 设为空数组
			if !errors.Is(err, badger.ErrKeyNotFound) { //如果真的发生异常，则开一个新事务继续读后面的 key
				txn.Discard()
				txn = s.db.NewTransaction(false)
			}
		}
	}
	txn.Discard() //只读事务调 Discard 就可以了，不需要调 Commit(实际上 Commit 内部也会调 Discard)
	return values, err
}

func (s *BadgerDB) Delete(k []byte) error {
	err := s.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(k)
	})
	return err
}

func (s *BadgerDB) BatchDelete(keys [][]byte) error {
	var err error
	txn := s.db.NewTransaction(true)
	for _, key := range keys {
		if err = txn.Delete(key); err != nil {
			_ = txn.Commit() //发生异常时就提交老事务，然后开一个新事务，重试 delete
			txn = s.db.NewTransaction(true)
			_ = txn.Delete(key)
		}
	}
	txn.Commit()
	return err
}

// IterateKey 只遍历 key。key 是全部存在 LSM tree 上的，只需要读内存，所以很快
func (s *BadgerDB) IterateKey(fn func(k []byte) error) int64 {
	var total int64
	_ = s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false //只需要读 key，所以把 PrefetchValues 设为 false
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			if err := fn(k); err == nil {
				atomic.AddInt64(&total, 1)
			}
		}
		return nil
	})
	return atomic.LoadInt64(&total)
}

// IterateDB 遍历整个 DB
func (s *BadgerDB) IterateDB(fn func(k, v []byte) error) int64 {
	var total int64
	_ = s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			key := item.Key()

			var ival []byte
			err := item.Value(func(val []byte) error {
				ival = val
				return nil
			})

			if err != nil {
				continue
			}
			if err := fn(key, ival); err == nil {
				atomic.AddInt64(&total, 1)
			}
		}
		return nil
	})
	return atomic.LoadInt64(&total)
}

// Has 判断某个 key 是否存在
func (s *BadgerDB) Has(k []byte) bool {
	//db.View 相当于打开了一个读写事务:db.NewTransaction(true)(用 db.Update 的好处在于不用显式调用 Txn.Discard())
	err := s.db.View(func(txn *badger.Txn) error {
		_, err := txn.Get(k)
		if err != nil {
			return err
		}
		return nil
	})
	// err 不为 nil 就假设 key 不存在
	if err != nil {
		return false
	}
	return true
}
