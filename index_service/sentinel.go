package index_service

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/hjrbill/quicker/gen"
	qlog "github.com/hjrbill/quicker/pkg/log"
	"github.com/hjrbill/quicker/pkg/service_hub"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"sync"
	"sync/atomic"
	"time"
)

var _ IIndexer = (*Sentinel)(nil)

// Sentinel 分布式 Work 的控制
type Sentinel struct {
	hub      service_hub.IHub // 服务注册中心
	connPool sync.Map         // 连接池（与各个 IndexServiceWorker 的连接的缓存）
}

// NewSentinel 默认使用 HubProxy，可使用 WithHub 设置自定义 Hub
// @etcdServers etcd 配置
// @loadBalancer 负载均衡算法
// @qps 限流 QPS
func NewSentinel(etcdServers []string, loadBalancer service_hub.LoadBalancer, qps int) *Sentinel {
	return &Sentinel{
		hub:      service_hub.GetServiceHubProxy(etcdServers, 10, loadBalancer, qps),
		connPool: sync.Map{},
	}
}

// WithHub 设置自定义的 Hub
func (s *Sentinel) WithHub(hub service_hub.IHub) *Sentinel {
	s.hub = hub
	return s
}

func (s *Sentinel) Close() (err error) {
	s.connPool.Range(func(key, value interface{}) bool {
		conn := value.(*grpc.ClientConn)
		closeErr := conn.Close()
		if closeErr != nil {
			// 不返回 false（中断）,继续处理剩下的连接，但报告错误提醒
			qlog.Errorf("关闭连接 %s 失败: %s", key, closeErr.Error())
			err = closeErr
		}
		s.connPool.Delete(key) // 从连接池中删除，便于用户的二次处理
		return true
	})
	s.hub.Close()
	return err
}

func (s *Sentinel) GetGrpcConn(endpoint string) *grpc.ClientConn {
	if v, ok := s.connPool.Load(endpoint); ok { // 从连接池中获取连接
		conn := v.(*grpc.ClientConn)
		if conn.GetState() == connectivity.Shutdown || conn.GetState() == connectivity.TransientFailure { // 如果连接已关闭或出于故障，则从连接池中删除
			s.connPool.Delete(endpoint)
			return nil
		} else {
			return conn
		}
	}
	// 连接池中没有连接，尝试建立连接
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	conn, err := grpc.DialContext(
		ctx,
		endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(), // 更改为同步连接，使得阻塞和超时控制生效
	)
	if err != nil {
		qlog.Warnf("连接 %s 失败: %s", endpoint, err.Error())
		return nil
	}
	s.connPool.Store(endpoint, conn)
	qlog.Infof("成功与 %s 建立连接", endpoint)
	return conn
}

// AddDoc 通过负载均衡算法获取可用的 IndexServiceWorker 并向其添加文档
func (s *Sentinel) AddDoc(doc pb.Document) (int, error) {
	endpoint := s.hub.GetEndpoint(INDEX_SERVICE)
	if endpoint == "" {
		return 0, errors.New("没有可用的 IndexServiceWorker")
	}
	conn := s.GetGrpcConn(endpoint)
	if conn == nil {
		return 0, fmt.Errorf("连接 IndexServiceWorker %s 失败", endpoint)
	}
	affected, err := pb.NewIndexServiceClient(conn).AddDocument(context.Background(), &doc)
	if err != nil {
		return 0, err
	}
	qlog.Infof("向 IndexServiceWorker %s 添加文档共 %d 条", endpoint, affected.Count)
	return int(affected.Count), nil
}

// DeleteDoc 向所有 IndexServiceWorker 广播，要求删除文档
func (s *Sentinel) DeleteDoc(id string) (int, error) {
	var affectedCount int32

	endpoints := s.hub.GetEndpoints(INDEX_SERVICE)
	if len(endpoints) == 0 {
		return 0, errors.New("没有可用的 IndexServiceWorker")
	}
	wg := sync.WaitGroup{}
	wg.Add(len(endpoints))
	for _, endpoint := range endpoints {
		go func(endpoint string) {
			defer wg.Done()
			conn := s.GetGrpcConn(endpoint)
			if conn != nil {
				affected, err := pb.NewIndexServiceClient(conn).DeleteDocument(context.Background(), &pb.ID{ID: id})
				if err != nil {
					qlog.Warnf("从 IndexServiceWorker %s 删除文档失败: %s", endpoint, err.Error())
				} else {
					if affected.Count > 0 {
						atomic.AddInt32(&affectedCount, affected.Count)
						qlog.Infof("从 IndexServiceWorker %s 删除文档共 %d 条", endpoint, affected.Count)
					}
				}
			}
		}(endpoint)
	}
	wg.Wait()
	return int(atomic.LoadInt32(&affectedCount)), nil
}

// Search 向所有 IndexServiceWorker 广播，进行检索，并合并结果
func (s *Sentinel) Search(query *pb.TermQuery, onFlag, offFlag uint64, orFlags []uint64) ([]*pb.Document, error) {
	endpoints := s.hub.GetEndpoints(INDEX_SERVICE)
	if len(endpoints) == 0 {
		return nil, errors.New("没有可用的 IndexServiceWorker")
	}

	docs := make([]*pb.Document, 0, len(endpoints)*10)
	result := make(chan *pb.Document, len(endpoints)*10)
	wg := sync.WaitGroup{}
	wg.Add(len(endpoints))
	for _, endpoint := range endpoints {
		go func(endpoint string) {
			defer wg.Done()
			conn := s.GetGrpcConn(endpoint)
			if conn != nil {
				docs, err := pb.NewIndexServiceClient(conn).Search(context.Background(), &pb.SearchRequest{
					TermQuery: query,
					OnFlag:    onFlag,
					OffFlag:   offFlag,
					OrFlag:    orFlags,
				})
				if err != nil {
					qlog.Warnf("从 IndexServiceWorker %s 搜索文档失败: %s", endpoint, err.Error())
				} else {
					if len(docs.Documents) > 0 {
						qlog.Infof("从 IndexServiceWorker %s 搜索得到文档共 %d 条", endpoint, len(docs.Documents))
						for _, doc := range docs.Documents {
							result <- doc
						}
					}
				}
			}
		}(endpoint)
	}

	resultFinish := make(chan struct{}, len(endpoints)*10)
	go func() {
		for {
			doc, ok := <-result
			if !ok {
				break
			}
			docs = append(docs, doc)
		}
		resultFinish <- struct{}{}
	}()
	wg.Wait()
	close(result)
	<-resultFinish
	return docs, nil
}

// Count 向所有 IndexServiceWorker 广播，要求获取文档总数
func (s *Sentinel) Count() int32 {
	var count int32

	endpoints := s.hub.GetEndpoints(INDEX_SERVICE)
	if len(endpoints) == 0 {
		return count
	}
	wg := sync.WaitGroup{}
	wg.Add(len(endpoints))
	for _, endpoint := range endpoints {
		go func(endpoint string) {
			defer wg.Done()
			conn := s.GetGrpcConn(endpoint)
			if conn != nil {
				affected, err := pb.NewIndexServiceClient(conn).Count(context.Background(), &pb.CountRequest{})
				if err != nil {
					qlog.Warnf("从 IndexServiceWorker %s 获取 doc count 失败: %s", endpoint, err.Error())
				} else {
					if affected.Count > 0 {
						atomic.AddInt32(&count, affected.Count)
						qlog.Infof("从 IndexServiceWorker %s 获取 doc count共 %d 条", endpoint, affected.Count)
					}
				}
			}
		}(endpoint)
	}
	wg.Wait()
	return atomic.LoadInt32(&count)
}
