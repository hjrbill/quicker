package test

import (
	"errors"
	"fmt"
	"log"
	"quicker/internal/kvdb"
	"testing"
)

// 本文件中没有进行实际的测试，只是为了方便测试，定义了测试该 interface 的所有功能所需的全局变量和单测函数

var (
	db     kvdb.IKeyValueDB
	setup  func()
	windup func()
)

func init() {
	windup = func() {
		err := db.Close()
		if err != nil {
			log.Printf("db close failed, err: %v", err)
		}
	}
}

func testSetAndGetAndDelete(db kvdb.IKeyValueDB) error {
	k1 := []byte("k1")
	v1 := []byte("v1")
	k2 := []byte("k2")
	v2 := []byte("v2")
	// 写入<k, v>
	err := db.Set(k1, v1)
	if err != nil {
		return err
	}
	err = db.Set(k2, v2)
	if err != nil {
		return err
	}
	// 读取<k, v>
	v, err := db.Get(k1)
	if err != nil {
		return err
	}
	fmt.Println("v1 =", string(v))
	v, err = db.Get(k2)
	if err != nil {
		return err
	}
	fmt.Println("v2 =", string(v))
	// 删除<k, v>
	err = db.Delete(k1)
	if err != nil {
		return err
	}
	err = db.Delete(k2)
	if err != nil {
		return err
	}
	// 读取<k, v>
	_, err = db.Get(k1)
	if err == nil {
		return errors.New("key 应被删除，却能读出数据")
	}
	v, err = db.Get(k2)
	if err == nil {
		return errors.New("key 应被删除，却能读出数据")
	}
	//判断 key 是否存在
	fmt.Printf("k1 存在：%t\n", db.Has(k1))
	fmt.Printf("k2 存在：%t\n", db.Has(k2))
	return nil
}

func testBatchSetAndBatchGetAndBatchDelete(db kvdb.IKeyValueDB) error {
	keys := [][]byte{[]byte("k1"), []byte("k2"), []byte("k3")}
	values := [][]byte{[]byte("v1"), []byte("v2"), []byte("v3")}

	// 批量写入
	err := db.BatchSet(keys, values)
	if err != nil {
		return err
	}
	// 批量读取
	vs, err := db.BatchGet(keys)
	if err != nil {
		return err
	}
	for i, v := range vs {
		fmt.Printf("v%d = %s\n", i, string(v))
	}
	// 批量删除
	err = db.BatchDelete(keys)
	if err != nil {
		return err
	}
	// 尝试读取被删除的数据（注意：不能用 db.BatchGet(),因为其获取空 key 时不会返回错误）
	_, err = db.Get(keys[0])
	if err == nil {
		return errors.New("k1 应被删除，却能读出数据")
	}
	_, err = db.Get(keys[1])
	if err == nil {
		return errors.New("k2 应被删除，却能读出数据")
	}
	_, err = db.Get(keys[2])
	if err == nil {
		return errors.New("k3 应被删除，却能读出数据")
	}
	// 判断 key 是否存在
	fmt.Printf("k1 存在 %t\n", db.Has(keys[0]))
	fmt.Printf("k2 存在 %t\n", db.Has(keys[1]))
	fmt.Printf("k3 存在 %t\n", db.Has(keys[2]))
	return nil
}

func testIterateDB(db kvdb.IKeyValueDB) error {
	keys := [][]byte{[]byte("k1"), []byte("k2"), []byte("k3")}
	values := [][]byte{[]byte("v1"), []byte("v2"), []byte("v3")}
	err := db.BatchSet(keys, values)
	if err != nil {
		return err
	}

	total := db.IterateDB(func(k, v []byte) error {
		fmt.Printf("key = %s, value = %s\n", k, v)
		return nil
	})
	fmt.Printf("total: %d\n", total)
	return nil
}

func testIterateKey(db kvdb.IKeyValueDB) error {
	keys := [][]byte{[]byte("k1"), []byte("k2"), []byte("k3")}
	values := [][]byte{[]byte("v1"), []byte("v2"), []byte("v3")}
	err := db.BatchSet(keys, values)
	if err != nil {
		return err
	}

	total := db.IterateKey(func(k []byte) error {
		fmt.Printf("key = %s\n", k)
		return nil
	})
	fmt.Printf("total: %d\n", total)
	return nil
}

// 总的 kv_db 的测试流水线
func testPipeline(t *testing.T) {
	// 初始化环境
	setup()
	defer windup()
	// 1. 测试单独写入，读取，删除
	err := testSetAndGetAndDelete(db)
	if err != nil {
		log.Printf("test failed: %v", err)
		t.Fail()
	}
	// 2. 测试批量写入，读取，删除
	err = testBatchSetAndBatchGetAndBatchDelete(db)
	if err != nil {
		log.Printf("test failed: %v", err)
		t.Fail()
	}
	// 3. 测试 IterateDB
	err = testIterateDB(db)
	if err != nil {
		log.Printf("test failed: %v", err)
		t.Fail()
	}
	// 4. 测试 IterateKey
	err = testIterateKey(db)
	if err != nil {
		log.Printf("test failed: %v", err)
		t.Fail()
	}
}
