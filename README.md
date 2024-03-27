# Quicker
一个简单、轻量级的分布式全文索引框架

### Translations:
- ### [English](/README_EN.md)

## 倒排索引

<img src="demo/views/img/倒排索引.png" width="500"/>    

- 倒排索引的 list 是有序的，便于多条倒排链快速求交集。
- DocId 是 quicker 系统内部给 doc 生成的自增 id，用于 SkipList 的排序。
- Id 是 doc 在业务侧的 ID。
- BitsFeature 是 uint64，可以把 doc 的属性编码成 bit 流，遍历倒排索引的同时完成部分筛选功能。

## 正排索引
暂时只支持 badger 和 bolt 两种数据库存储 doc 的详情。

## 分布式索引
各索引服务器之间通过 grpc 通信，通过 etcd 实现服务注册与发现。

## 使用方法
`go get -u github.com/hjrbill/quicker`

## Quickstart
### 假设业务 doc 为 Book
```go 
type Book struct {
	ISBN    string
	Title   string
	Author  string
	Price   float64
	Content string
}
// 业务侧自行实现 doc 的序列化和反序列化
func (book *Book) Serialize() []byte {
}
func DeserializeBook(v []byte) *Book {
}
```
### 初始化 quicker
```go
dbType := kvdb.BADGER   //或 kvdb.BOLT
dbPath := "data/local_db/book_badger"   //正排索引的存储路径
docNum := 10000    //预估索引里将存储多少文档
quicker := new(index_service.Indexer)
if err := quicker.Init(docNum, dbType, dbPath); err != nil {
    panic(err)
}
defer quicker.Close()
```
### 添加 doc
```go
book := Book{}
doc := types.Document{
		Id:          book.ISBN,
		BitsFeature: 0b10011, //二进制
		Keywords:    []*types.Keyword{{Field: "content", Word: "唐朝"}, {Field: "content", Word: "文物"}, {Field: "title", Word: book.Title}},
		Bytes:       book.Serialize(),
	}
quicker.AddDoc(doc)
```
### 删除 doc
```go 
quicker.DeleteDoc(doc.Id)
```
### 检索
```go 
q1 := types.NewTermQuery("title", "生命起源")
q2 := types.NewTermQuery("content", "文物")
q3 := types.NewTermQuery("title", "中国历史")
q4 := types.NewTermQuery("content", "文物")
q5 := types.NewTermQuery("content", "唐朝")

// 支持任意复杂的 And 和 Or 的组合。And 要求同时命中，Or 只要求命中一个
query := (q1.Or(q2)).And((q3.Or(q4)).And(q5))
var onFlag uint64 = 0b10000    //要求 doc.BitsFeature 的对应位必须都是 1
var offFlag uint64 = 0b01000    //要求 doc.BitsFeature 的对应位必须都是 0
orFlags := []uint64{uint64(0b00010), uint64(0b00101)}    //要求 doc.BitsFeature 的对应位至少有一个是 1
docs := quicker.Search(query, onFlag, offFlag, orFlags) //检索
for _, doc := range docs {
    book := DeserializeBook(doc.Bytes) //检索的结果是二进流，需要自反序列化
    if book != nil {
        fmt.Printf("%s %s %s %.1f\n", book.ISBN, book.Title, book.Author, book.Price)
    }
}
```