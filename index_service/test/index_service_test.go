package test

import (
	"context"
	"fmt"
	"github.com/hjrbill/quicker/index_service"
	"github.com/hjrbill/quicker/internal/kvdb"
	"github.com/hjrbill/quicker/pb"
	qlog "github.com/hjrbill/quicker/pkg/log"
	"github.com/hjrbill/quicker/pkg/util"
	"net"
	"strconv"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	servicePort = 5678
)

// 启动 grpc server
func StartService() {
	// 监听本地的 5678 端口
	lis, err := net.Listen("tcp", "127.0.0.1:"+strconv.Itoa(servicePort))
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()
	service := new(index_service.IndexServiceWork)
	err = service.Init(50000, kvdb.BADGER, util.RootPath+"temp/book_badger")
	if err != nil {
		panic(err)
	} //不进行服务注册，client 直连 server
	service.Indexer.LoadFromForwardIndexFile() //从文件中加载索引数据
	// 注册服务的具体实现
	pb.RegisterIndexServiceServer(server, service)
	go func() {
		// 启动服务
		fmt.Printf("start grpc server on port %d\n", servicePort)
		err = server.Serve(lis) //Serve 会一直阻塞，所以放到一个协程里异步执行
		if err != nil {
			panic(err)
		}
	}()
}

func TestIndexService(t *testing.T) {
	StartService()              //server 和 client 分到不同的协程里去。实际中，server 和 client 是部署在不同的机器上
	time.Sleep(1 * time.Second) //等 server 启动完毕

	//连接到服务端
	conn, err := grpc.DialContext(
		context.Background(),
		"127.0.0.1:"+strconv.Itoa(servicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()), //Credential 即使为空，也必须设置
	)
	if err != nil {
		fmt.Printf("dial failed: %s", err)
		return
	}
	//创建 client
	client := pb.NewIndexServiceClient(conn)

	//测试 Search 接口
	query := pb.NewTermQuery("content", "文物")
	query = query.And(pb.NewTermQuery("content", "唐朝"))
	request := &pb.SearchRequest{
		TermQuery: query,
	}
	resp, err := client.Search(context.Background(), request)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	} else {
		docId := ""
		qlog.Debugf("%d documents were found", len(resp.Documents))
		for _, doc := range resp.Documents {
			book := DeserializeBook(doc.Bytes) //检索的结果是二进流，需要自反序列化
			if book != nil {
				fmt.Printf("%s %s %s %s %.1f\n", doc.Id, book.ISBN, book.Title, book.Author, book.Price)
				docId = doc.Id
			}
		}
		//测试 Delete 接口
		if len(docId) > 0 {
			affect, err := client.DeleteDocument(context.Background(), &pb.ID{ID: docId})
			if err != nil {
				fmt.Println(err)
				t.Fail()
			} else {
				fmt.Printf("删除%d个 doc\n", affect.Count)
			}
		}
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
		affect, err := client.AddDocument(context.Background(), &doc)
		if err != nil {
			fmt.Println(err)
			t.Fail()
		} else {
			fmt.Printf("添加%d个 doc\n", affect.Count)
		}
		//测试 Search 接口
		request := &pb.SearchRequest{
			TermQuery: query,
		}
		result, err := client.Search(context.Background(), request)
		if err != nil {
			fmt.Println(err)
			t.Fail()
		} else {
			for _, doc := range result.Documents {
				book := DeserializeBook(doc.Bytes) //检索的结果是二进流，需要自反序列化
				if book != nil {
					fmt.Printf("%s %s %s %s %.1f\n", doc.Id, book.ISBN, book.Title, book.Author, book.Price)
				}
			}
		}
	}
}
