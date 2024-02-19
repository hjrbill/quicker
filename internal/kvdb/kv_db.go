package kvdb

import (
	"errors"
	"log"
	"os"
	"strings"
)

var ErrDataNotFound = errors.New("data not found")

type DBType int

const (
	BOLT   DBType = iota // go 编写，基于 B+树
	BADGER               // c++编写，基于 LSM
	REDIS                // ANSI C 编写，基于 内存存储
)

type IKeyValueDB interface {
	GetDBPath() string                                // 获取数据存储的目录
	Open() error                                      // 建立 DB 连接
	Close() error                                     // 关闭连接
	Set(key, value []byte) error                      // 写入 key-value
	BatchSet(keys, values [][]byte) error             // 批量写入
	Get(key []byte) ([]byte, error)                   // 根据 key 获取 value
	BatchGet(keys [][]byte) ([][]byte, error)         // 批量获取 value
	Delete(key []byte) error                          // 删除 key-value
	BatchDelete(keys [][]byte) error                  // 批量删除
	IterateKey(fn func(key []byte) error) int64       // 遍历所有 key，返回获得的条数
	IterateDB(fn func(key, value []byte) error) int64 // 遍历整个数据库，返回获得的条数
	Has(key []byte) bool                              // 判断 key 是否存在
}

func NewKVDB(dbType DBType, path string) (IKeyValueDB, error) {
	paths := strings.Split(path, "/")
	parentPath := strings.Join(paths[0:len(paths)-1], "/") // 去除路径最后一级以获取父目录

	info, err := os.Stat(parentPath)
	if os.IsNotExist(err) {
		// TODO：log 应该被替换
		log.Printf("create dir: %s", parentPath)
		err := os.MkdirAll(parentPath, os.ModePerm)
		if err != nil {
			return nil, errors.New("create dir error: " + err.Error())
		}
	} else {
		// 检查父目录是否为目录
		if !info.IsDir() {
			return nil, errors.New("parent path was exist and it is not dir")
		}
		//// 检查父目录是否为普通文件
		//if info.Mode().IsRegular() { //如果父路径是个普通文件，则把它删掉
		//	log.Printf("%s is a regular file, will delete it", parentPath)
		//	os.Remove(parentPath)
		//}
	}
	var db IKeyValueDB
	switch dbType {
	case BADGER:
		db = new(BadgerDB).WithPath(path)
	case REDIS:
		// TODO:暂时使用 bolt 代替，补全对 Redis 的支持
		db = new(BoltDB).WithPath(path).WithBucket("default")
	default:
		// 默认使用 bolt
		db = new(BoltDB).WithPath(path).WithBucket("default")
	}
	err = db.Open()
	if err != nil {
		return nil, err
	}
	return db, nil
}
