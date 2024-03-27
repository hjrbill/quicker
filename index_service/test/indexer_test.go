package test

import (
	"bytes"
	"encoding/gob"
	"fmt"
	pb "github.com/hjrbill/quicker/gen"
	"github.com/hjrbill/quicker/index_service"
	"github.com/hjrbill/quicker/internal/kvdb"
	qlog "github.com/hjrbill/quicker/pkg/log"
	"github.com/hjrbill/quicker/pkg/path"
	"strings"
	"testing"
)

type Book struct {
	ISBN    string
	Title   string
	Author  string
	Price   float64
	Content string
}

// Serialize 图书序列化。序列化和反序列化由调用方决定，这不是 quicker 负责的范畴。
func (book *Book) Serialize() []byte {
	var value bytes.Buffer
	encoder := gob.NewEncoder(&value) //gob 是 go 自带的序列化方法，当然也可以用 protobuf 等其它方式
	err := encoder.Encode(book)
	if err == nil {
		return value.Bytes()
	} else {
		fmt.Println("序列化图书失败", err)
		return []byte{}
	}
}

// DeserializeBook  图书反序列化
func DeserializeBook(v []byte) *Book {
	buf := bytes.NewReader(v)
	dec := gob.NewDecoder(buf)
	var data = Book{}
	err := dec.Decode(&data)
	if err == nil {
		return &data
	} else {
		fmt.Println("反序列化图书失败", err)
		return nil
	}
}

var (
	dbType = kvdb.BADGER
	dbPath = path.RootPath + "temp/test/badger_db"
)

func TestSearch(t *testing.T) {
	quicker := new(index_service.Indexer)
	if err := quicker.Init(100, dbType, dbPath); err != nil {
		fmt.Println(err)
		t.Fail()
	}
	defer quicker.Close()

	book1 := Book{
		ISBN:    "315246546",
		Title:   "计算机系列丛书",
		Author:  "张三",
		Price:   59.0,
		Content: "冰雪奇缘 2 中文版电影原声带 (Frozen 2 (Mandarin Original Motion Picture",
	}
	book2 := Book{
		ISBN:    "436246383",
		Title:   "中国历史",
		Author:  "李四",
		Price:   39.0,
		Content: "冰雪奇缘 2 中文版电影原声带 (Frozen 2 (Mandarin Original Motion Picture",
	}
	book3 := Book{
		ISBN:    "54325435634",
		Title:   "生命起源",
		Author:  "赵六",
		Price:   49.0,
		Content: "冰雪奇缘 2 中文版电影原声带 (Frozen 2 (Mandarin Original Motion Picture",
	}

	doc1 := pb.Document{
		Id:          book1.ISBN,
		BitsFeature: 0b10101, //二进制
		Keywords:    []*pb.Keyword{{Field: "content", Word: "机器学习"}, {Field: "content", Word: "神经网络"}, {Field: "title", Word: book1.Title}},
		Bytes:       book1.Serialize(), //写入索引时需要自行序列化
	}
	doc2 := pb.Document{
		Id:          book2.ISBN,
		BitsFeature: 0b10011, //二进制
		Keywords:    []*pb.Keyword{{Field: "content", Word: "唐朝"}, {Field: "content", Word: "文物"}, {Field: "title", Word: book2.Title}},
		Bytes:       book2.Serialize(),
	}
	doc3 := pb.Document{
		Id:          book3.ISBN,
		BitsFeature: 0b11101, //二进制
		Keywords:    []*pb.Keyword{{Field: "content", Word: "动物"}, {Field: "content", Word: "文物"}, {Field: "title", Word: book3.Title}},
		Bytes:       book3.Serialize(),
	}

	cnt1, err := quicker.AddDoc(doc1)
	if err != nil {
		t.Fail()
	}
	cnt2, err := quicker.AddDoc(doc2)
	if err != nil {
		t.Fail()
	}
	cnt3, err := quicker.AddDoc(doc3)
	if err != nil {
		t.Fail()
	}
	qlog.Debugf("Successfully added %d documents", cnt1+cnt2+cnt3)

	q1 := pb.NewTermQuery("title", "生命起源")
	q2 := pb.NewTermQuery("content", "文物")
	q3 := pb.NewTermQuery("title", "中国历史")
	q4 := pb.NewTermQuery("content", "文物")
	q5 := pb.NewTermQuery("content", "唐朝")

	q6 := q1.And(q2)
	q7 := q3.And(q4).And(q5)

	q8 := q6.Or(q7)

	var onFlag uint64 = 0b10000
	var offFlag uint64 = 0b01000
	orFlags := []uint64{uint64(0b00010), uint64(0b00101)}
	docs, err := quicker.Search(q8, onFlag, offFlag, orFlags) //检索
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	qlog.Debugf("%d documents were found", len(docs))

	for _, doc := range docs {
		book := DeserializeBook(doc.Bytes) //检索的结果是二进流，需要自反序列化
		if book != nil {
			fmt.Printf("%s %s %s %.1f\n", book.ISBN, book.Title, book.Author, book.Price)
		}
	}
	fmt.Println(strings.Repeat("-", 50))

	_, err = quicker.DeleteDoc(doc2.Id)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	docs, err = quicker.Search(q8, onFlag, offFlag, orFlags) //检索
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	qlog.Debugf("%d documents were found", len(docs))
	for _, doc := range docs {
		book := DeserializeBook(doc.Bytes) //检索的结果是二进流，需要自反序列化
		if book != nil {
			fmt.Printf("%s %s %s %.1f\n", book.ISBN, book.Title, book.Author, book.Price)
		}
	}
	fmt.Println(strings.Repeat("-", 50))

	_, err = quicker.AddDoc(doc2)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	docs, err = quicker.Search(q8, onFlag, offFlag, orFlags) //检索
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	qlog.Debugf("%d documents were found", len(docs))
	for _, doc := range docs {
		book := DeserializeBook(doc.Bytes) //检索的结果是二进流，需要自反序列化
		if book != nil {
			fmt.Printf("%s %s %s %.1f\n", book.ISBN, book.Title, book.Author, book.Price)
		}
	}
	fmt.Println(strings.Repeat("-", 50))
}

func TestLoadFromIndexFile(t *testing.T) {
	indexer := new(index_service.Indexer)
	if err := indexer.Init(100, dbType, dbPath); err != nil {
		fmt.Println(err)
		t.Fail()
		return
	}
	defer indexer.Close()

	n := indexer.LoadFromForwardIndexFile()
	if n == 0 {
		return
	} else {
		qlog.Debugf("load %d document from invere document", n)
	}

	q1 := pb.NewTermQuery("title", "生命起源")
	q2 := pb.NewTermQuery("content", "文物")
	q3 := pb.NewTermQuery("title", "中国历史")
	q4 := pb.NewTermQuery("content", "文物")
	q5 := pb.NewTermQuery("content", "唐朝")

	q6 := q1.And(q2)
	q7 := q3.And(q4).And(q5)

	q8 := q6.Or(q7)

	var onFlag uint64 = 0b10000
	var offFlag uint64 = 0b01000
	orFlags := []uint64{uint64(0b00010), uint64(0b00101)}
	docs, err := indexer.Search(q8, onFlag, offFlag, orFlags) //检索
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	qlog.Debugf("%d documents were found", len(docs))
	for _, doc := range docs {
		book := DeserializeBook(doc.Bytes) //检索的结果是二进流，需要自反序列化
		if book != nil {
			fmt.Printf("%s %s %s %.1f\n", book.ISBN, book.Title, book.Author, book.Price)
		}
	}
	fmt.Println(strings.Repeat("-", 50))

	book2 := Book{
		ISBN:    "436246383",
		Title:   "中国历史",
		Author:  "李四",
		Price:   39.0,
		Content: "冰雪奇缘 2 中文版电影原声带 (Frozen 2 (Mandarin Original Motion Picture",
	}
	doc2 := pb.Document{
		Id:          book2.ISBN,
		BitsFeature: 0b10011, //二进制
		Keywords:    []*pb.Keyword{{Field: "content", Word: "唐朝"}, {Field: "content", Word: "文物"}, {Field: "title", Word: book2.Title}},
		Bytes:       book2.Serialize(),
	}

	_, err = indexer.DeleteDoc(doc2.Id)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	docs, err = indexer.Search(q8, onFlag, offFlag, orFlags) //检索
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	qlog.Debugf("%d documents were found", len(docs))
	for _, doc := range docs {
		book := DeserializeBook(doc.Bytes) //检索的结果是二进流，需要自反序列化
		if book != nil {
			fmt.Printf("%s %s %s %.1f\n", book.ISBN, book.Title, book.Author, book.Price)
		}
	}
	fmt.Println(strings.Repeat("-", 50))

	_, err = indexer.AddDoc(doc2)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	docs, err = indexer.Search(q8, onFlag, offFlag, orFlags) //检索
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	qlog.Debugf("%d documents were found", len(docs))
	for _, doc := range docs {
		book := DeserializeBook(doc.Bytes) //检索的结果是二进流，需要自反序列化
		if book != nil {
			fmt.Printf("%s %s %s %.1f\n", book.ISBN, book.Title, book.Author, book.Price)
		}
	}
	fmt.Println(strings.Repeat("-", 50))
}
