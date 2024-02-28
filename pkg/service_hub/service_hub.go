package service_hub

import (
	"context"
	"errors"
	qlog "github.com/hjrbill/quicker/pkg/log"
	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	etcd "go.etcd.io/etcd/client/v3"
	"strings"
	"sync"
	"time"
)

var _ IHub = (*ServiceHub)(nil)

type ServiceHub struct {
	client       *etcd.Client
	watched      sync.Map // 存储对服务的监听，key 为服务名，value 为 etcd.WatchChan
	heartbeat    int64    // 心跳频率（续约周期）
	loadBalancer LoadBalancer
}

var (
	serviceHub     *ServiceHub
	serviceHubOnce sync.Once
)

// GetServiceHub
// @param etcdServers etcd 集群地址
// @param heartbeat 心跳频率（续约周期）
// @param loadBalancer 负载均衡策略
func GetServiceHub(etcdServers []string, heartbeat int64, loadBalancer LoadBalancer) *ServiceHub {
	if serviceHub == nil {
		serviceHubOnce.Do(func() {
			client, err := etcd.New(etcd.Config{
				Endpoints:   etcdServers,
				DialTimeout: 3 * time.Second,
				Context:     context.Background(),
			})
			if err != nil {
				qlog.Fatalf("连接 etcd 失败: %v", err)
			}

			serviceHub = &ServiceHub{
				client:       client,
				heartbeat:    heartbeat,
				loadBalancer: loadBalancer,
			}
		})
	}
	return serviceHub
}

func (s *ServiceHub) Close() {
	if s.client != nil {
		err := s.client.Close()
		if err != nil {
			qlog.Errorf("关闭 etcd 客户端失败: %v", err)
		}
	}
}

func (s *ServiceHub) Register(serviceName string, endpoint string, leaseID etcd.LeaseID) (etcd.LeaseID, error) {
	if s.client == nil {
		return 0, errors.New("etcd 客户端未初始化")
	}

	ctx := context.Background()
	if leaseID == 0 {
		// 先获取一份租约
		if lease, err := s.client.Grant(ctx, s.heartbeat); err != nil {
			qlog.Warnf("获取租约失败: %v", err)
			return 0, err
		} else {
			prefix := strings.TrimRight(SERVICE_ROOT_PATH, "/") + "/" + serviceName + "/" + endpoint
			// 通过租约将节点写入
			if _, err := s.client.Put(ctx, prefix, endpoint, etcd.WithLease(lease.ID)); err != nil {
				qlog.Warnf("将节点 %s 写入 %s 服务失败: %v", endpoint, serviceName, err)
				return 0, err
			}
			qlog.Infof("将节点 %s 写入 %s 服务成功", endpoint, serviceName)
			return lease.ID, nil
		}
	} else {
		_, err := s.client.KeepAliveOnce(ctx, leaseID) // 尝试发送一次续约
		if errors.Is(err, rpctypes.ErrLeaseNotFound) { // 如果租约不存在，走注册流程
			return s.Register(serviceName, endpoint, 0)
		} else if err != nil {
			qlog.Warnf("节点%s续约失败: %v", endpoint, err)
			return 0, err
		}
		return leaseID, nil
	}
}

func (s *ServiceHub) UnRegister(serviceName string, endpoint string) error {
	if s.client == nil {
		return errors.New("etcd 客户端未初始化")
	}

	prefix := strings.TrimRight(SERVICE_ROOT_PATH, "/") + "/" + serviceName + "/" + endpoint
	_, err := s.client.Delete(context.Background(), prefix)
	if err != nil {
		qlog.Warnf("从 %s 服务中注销节点 %s 失败: %v", serviceName, endpoint, err)
		return err
	}
	qlog.Infof("从 %s 服务中注销节点 %s 成功", serviceName, endpoint)
	return nil
}

// GetEndpoints 获取服务节点列表（自行实现节点选择）
func (s *ServiceHub) GetEndpoints(serviceName string) []string {
	prefix := strings.TrimRight(SERVICE_ROOT_PATH, "/") + "/" + serviceName
	resp, err := s.client.Get(context.Background(), prefix, etcd.WithPrefix()) // 尝试以服务名为前缀获取节点
	if err != nil {
		qlog.Warnf("获取 %s 服务节点失败: %v", serviceName, err)
		return nil
	}

	endpoints := make([]string, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		endpoints = append(endpoints, string(kv.Value))
	}
	return endpoints
}

// GetEndpoint 获取服务节点（可调用我方提供负载均衡或自行实现）
func (s *ServiceHub) GetEndpoint(serviceName string) string {
	if s.client == nil || s.loadBalancer == nil {
		return ""
	} else {
		return s.loadBalancer.Take(s.GetEndpoints(serviceName)) // 通过负载均衡算法从可用节点中获取一个
	}
}
