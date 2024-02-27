package service_hub

import (
	"context"
	etcd "go.etcd.io/etcd/client/v3"
	"golang.org/x/time/rate"
	qlog "quicker/pkg/log"
	"strings"
	"sync"
	"time"
)

var _ IHub = (*HubProxy)(nil)

// HubProxy ServiceHub 的代理，额外提供了限流与缓存的功能
type HubProxy struct {
	*ServiceHub
	endpointCache sync.Map      // 缓存服务节点地址
	limiter       *rate.Limiter // 令牌桶的限流器
}

var (
	proxy     *HubProxy
	proxyOnce sync.Once
)

// GetServiceHubProxy
// @param etcdServers etcd 集群地址
// @param heartbeat 心跳频率（续约周期）
// @param loadBalancer 负载均衡策略
// @param qps 限流频率（每秒产生的令牌数量）
func GetServiceHubProxy(etcdServers []string, heartbeat int64, loadBalancer LoadBalancer, qps int) *HubProxy {
	if proxy == nil {
		proxyOnce.Do(func() {
			serviceHub := GetServiceHub(etcdServers, heartbeat, loadBalancer)
			proxy = &HubProxy{
				ServiceHub: serviceHub,
				limiter:    rate.NewLimiter(rate.Every(time.Duration(1e9/qps)*time.Nanosecond), qps), // 每1e9/qps纳秒产生一个令牌（每秒共qps个）,桶大小为qps
			}
		})
	}
	return proxy
}

func (p *HubProxy) watchService(serviceName string) {
	if _, ok := p.watched.LoadOrStore(serviceName, struct{}{}); ok {
		return // 该 service 已经被监听，直接返回
	}

	prefix := strings.TrimRight(SERVICE_ROOT_PATH, "/") + "/" + serviceName
	watchChan := p.client.Watch(context.Background(), prefix, etcd.WithPrefix())
	go func() {
		for resp := range watchChan {
			for _, event := range resp.Events { // 获取到的是事件的合集
				qlog.Infof("监听到服务 %s 发生了变化：%s", serviceName, event.Type)
				// TODO:要不要根据 event.Type 判断是 put 还是 delete，以此减少更新的量
				// 如果发生了变化，则对对应服务进行一次全量更新
				path := strings.Split(string(event.Kv.Key), "/")
				if len(path) > 2 {
					serviceName := path[len(path)-2]
					endpoints := p.ServiceHub.GetEndpoints(serviceName)
					if len(endpoints) > 0 {
						p.endpointCache.Store(serviceName, endpoints) // 对于查询结果进行缓存
					} else {
						qlog.Infof("现已停止支持 %s 服务", serviceName)
						p.endpointCache.Delete(serviceName) // 该服务已无节点，删除缓存
					}
				}
			}
		}
	}()
}

// GetEndpoints 获取服务节点列表（自行实现节点选择），增加了限流保护与缓存
// 重写了 ServiceHub 中的 GetEndpoints，并且因为 GetEndpoint 调用了 GetEndpoints，所以实际上也改变了 GetEndpoint 的逻辑
func (p *HubProxy) GetEndpoints(serviceName string) []string {
	//if p.limiter.Allow() { // 如果未能获取令牌，视为无可用节点，直接返回空
	//	return nil
	//}
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	if p.limiter.Wait(ctx) != nil { // 如果 100 纳秒都没有令牌，视为存在网络错误或争抢严重，无可用节点，直接返回空
		return nil
	}

	p.watchService(serviceName)
	if endpoints, ok := p.endpointCache.Load(serviceName); ok {
		return endpoints.([]string)
	} else {
		endpoints := p.ServiceHub.GetEndpoints(serviceName) // 尝试获取服务节点
		if len(endpoints) > 0 {
			qlog.Infof("现以支持服务 %s", serviceName)          // 只有一个服务第一次被启用时，才可能到达此处，所以进行一次广播
			p.endpointCache.Store(serviceName, endpoints) // 对于查询结果进行缓存
		}
		return endpoints
	}
}
