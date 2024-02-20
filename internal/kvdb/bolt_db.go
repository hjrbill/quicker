package kvdb

import (
	"errors"
	"go.etcd.io/bbolt"
	"sync/atomic"
)

var _ = IKeyValueDB(&BoltDB{})

type BoltDB struct {
	db     *bbolt.DB
	path   string
	bucket []byte
}

func (b *BoltDB) WithPath(path string) *BoltDB {
	b.path = path
	return b
}

func (b *BoltDB) WithBucket(bucket string) *BoltDB {
	b.bucket = []byte(bucket)
	return b
}

func (b *BoltDB) GetDBPath() string {
	return b.path
}

func (b *BoltDB) Open() error {
	dir := b.GetDBPath()
	db, err := bbolt.Open(dir, 0o600, bbolt.DefaultOptions)
	if err != nil {
		return err
	}
	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(b.bucket)
		if err != nil {
			db.Close()
			return err
		}
		return err
	})
	if err != nil {
		db.Close()
		return err
	}
	b.db = db
	return nil
}

func (b *BoltDB) Close() error {
	err := b.db.Close()
	if err != nil {
		return err
	}
	return nil
}

func (b *BoltDB) Set(key, value []byte) error {
	return b.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(b.bucket)
		return bucket.Put(key, value)
	})
}

func (b *BoltDB) BatchSet(keys, values [][]byte) error {
	if len(keys) != len(values) {
		return errors.New("keys and values should be same length")
	}
	return b.db.Update(func(tx *bbolt.Tx) error {
		for i := range keys {
			err := tx.Bucket(b.bucket).Put(keys[i], values[i])
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// Get 在只读事务的上下文中执行 get
func (b *BoltDB) Get(k []byte) ([]byte, error) {
	var ival []byte
	err := b.db.View(func(tx *bbolt.Tx) error {
		ival = tx.Bucket(b.bucket).Get(k)
		return nil
	})
	if len(ival) == 0 {
		return nil, ErrDataNotFound
	}
	return ival, err
}

func (b *BoltDB) BatchGet(keys [][]byte) ([][]byte, error) {
	values := make([][]byte, len(keys))
	err := b.db.Batch(func(tx *bbolt.Tx) error {
		for i, key := range keys {
			ival := tx.Bucket(b.bucket).Get(key)
			values[i] = ival
		}
		return nil
	})
	return values, err
}

func (b *BoltDB) Delete(k []byte) error {
	return b.db.Update(func(tx *bbolt.Tx) error {
		return tx.Bucket(b.bucket).Delete(k)
	})
}

func (b *BoltDB) BatchDelete(keys [][]byte) error {
	err := b.db.Batch(func(tx *bbolt.Tx) error {
		for _, key := range keys {
			err := tx.Bucket(b.bucket).Delete(key)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (b *BoltDB) IterateKey(fn func(k []byte) error) int64 {
	var total int64
	_ = b.db.View(func(tx *bbolt.Tx) error {
		cur := tx.Bucket(b.bucket).Cursor()
		for k, _ := cur.First(); k != nil; k, _ = cur.Next() {
			err := fn(k)
			if err != nil {
				return err
			} else {
				atomic.AddInt64(&total, 1)
			}
		}
		return nil
	})
	return atomic.LoadInt64(&total)
}

func (b *BoltDB) IterateDB(fn func(k, v []byte) error) int64 {
	var total int64
	_ = b.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(b.bucket)
		c := bucket.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			err := fn(k, v)
			if err != nil {
				return err
			} else {
				atomic.AddInt64(&total, 1)
			}
		}
		return nil
	})
	return atomic.LoadInt64(&total)
}

func (b *BoltDB) Has(k []byte) bool {
	var value []byte
	err := b.db.View(func(tx *bbolt.Tx) error {
		value = tx.Bucket(b.bucket).Get(k)
		return nil
	})
	if err != nil || string(value) == "" {
		return false
	}
	return true
}
