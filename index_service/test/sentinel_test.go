package test

import (
	"fmt"
	pb "github.com/hjrbill/quicker/gen"
	"github.com/hjrbill/quicker/index_service"
	"github.com/hjrbill/quicker/internal/kvdb"
	qlog "github.com/hjrbill/quicker/pkg/log"
	"github.com/hjrbill/quicker/pkg/path"
	"github.com/hjrbill/quicker/pkg/service_hub"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"testing"
	"time"
)

var (
	//在一台机器上启多个 worker，模拟实际生产中多台机器上启动多个 worker
	etcdServers = []string{"127.0.0.1:2379"}
	workPorts   = []int{5678, 5679, 5660}

	workers []*index_service.IndexServiceWorker
)

func StartWorkers() {
	workers = make([]*index_service.IndexServiceWorker, 0, len(workPorts))
	for i, port := range workPorts {
		// 监听本地端口
		lis, err := net.Listen("tcp", "127.0.0.1:"+strconv.Itoa(port))
		if err != nil {
			panic(err)
		}

		server := grpc.NewServer()
		service := new(index_service.IndexServiceWorker)
		err = service.Init(50000, kvdb.BADGER, path.RootPath+"temp/book_badger_"+strconv.Itoa(i))
		if err != nil {
			panic(err)
		}
		service.Indexer.LoadFromForwardIndexFile() //从文件中加载索引数据
		// 注册服务的具体实现
		pb.RegisterIndexServiceServer(server, service)
		err = service.Register(port, etcdServers, &service_hub.RoundRobin{})
		if err != nil {
			panic(err)
		}
		go func(port int) {
			// 启动服务
			qlog.Debugf("start grpc server on port %d\n", port)
			err = server.Serve(lis) //Serve 会一直阻塞，所以放到一个协程里异步执行
			if err != nil {
				err := service.Close()
				if err != nil {
					qlog.Debugf("关闭 etcd 客户端失败: %v", err)
				}
				qlog.Debugf("start grpc server on port %d failed: %s\n", port, err)
			} else {
				workers = append(workers, service)
			}
		}(port)
	}
}

func StopWorkers() {
	for _, worker := range workers {
		err := worker.Close()
		if err != nil {
			qlog.Debugf("关闭 etcd 客户端失败: %v", err)
		}
	}
}

func TestIndexCluster(t *testing.T) {
	StartWorkers()
	time.Sleep(3 * time.Second) //等所有 worker 都启动完毕
	defer StopWorkers()

	sentinel := index_service.NewSentinel(etcdServers, &service_hub.RoundRobin{}, 100)
	//测试 Add 接口
	book := Book{
		ISBN:    "436246383",
		Title:   "上下五千年",
		Author:  "李四",
		Price:   39.0,
		Content: "冰雪奇缘 2 中文版电影原声带 (Frozen 2 (Mandarin Original Motion Picture",
	}
	doc := pb.Document{
		Id:          book.ISBN,
		BitsFeature: 0b10011, //二进制
		Keywords:    []*pb.Keyword{{Field: "content", Word: "唐朝"}, {Field: "content", Word: "文物"}, {Field: "title", Word: book.Title}},
		Bytes:       book.Serialize(),
	}

	n, err := sentinel.AddDoc(doc)
	if err != nil {
		qlog.Debugf("添加文档失败: %v", err)
		t.Fail()
	} else {
		qlog.Debugf("添加%d个 doc\n", n)
	}
	//测试 Search 接口
	query := pb.NewTermQuery("content", "文物")
	query = query.And(pb.NewTermQuery("content", "唐朝"))
	docs, err := sentinel.Search(query, 0, 0, nil)
	if err != nil {
		qlog.Debugf("检索文档失败: %v", err)
		t.Fail()
	}

	docId := ""
	if len(docs) == 0 {
		qlog.Debug("无搜索结果")
	} else {
		for _, doc := range docs {
			book := DeserializeBook(doc.Bytes) //检索的结果是二进流，需要自反序列化
			if book != nil {
				qlog.Debugf("%s %s %s %s %.1f\n", doc.Id, book.ISBN, book.Title, book.Author, book.Price)
				docId = doc.Id
			}
		}
	}

	//测试 Delete 接口
	if len(docId) > 0 {
		n, err := sentinel.DeleteDoc(docId)
		if err != nil {
			qlog.Debugf("删除文档失败: %v", err)
			t.Fail()
		}
		qlog.Debugf("删除%d个 doc\n", n)
	}

	//测试 Search 接口
	docs, err = sentinel.Search(query, 0, 0, nil)
	if err != nil {
		qlog.Debugf("检索文档失败: %v", err)
		t.Fail()
	}
	if len(docs) == 0 {
		qlog.Debug("无搜索结果")
	} else {
		for _, doc := range docs {
			book := DeserializeBook(doc.Bytes) //检索的结果是二进流，需要自反序列化
			if book != nil {
				qlog.Debugf("%s %s %s %s %.1f\n", doc.Id, book.ISBN, book.Title, book.Author, book.Price)
			}
			fmt.Println("book:", book)
		}
	}

	println("end ................")
	time.Sleep(3 * time.Second)
}
