package service_hub

import (
	"math/rand"
	"sync/atomic"
)

var (
	_ LoadBalancer = (*Random)(nil)
	_ LoadBalancer = (*RoundRobin)(nil)
)

// LoadBalancer 一个策略，定义了负载均衡的接口，用户可以自行实现，我方也提供了两种简单的负载均衡算法
type LoadBalancer interface {
	Take(endpoints []string) string // 如何从大量可用地址中通过负载均衡返回某地址
}

// Random 随机法
type Random struct {
}

func (r *Random) Take(endpoints []string) string {
	if endpoints == nil || len(endpoints) == 0 {
		return ""
	}
	index := rand.Intn(len(endpoints)) // 获取一个随机值
	return endpoints[index]
}

// RoundRobin 轮训法
type RoundRobin struct {
	cnt uint64 // 记录调用次数
}

func (r *RoundRobin) Take(endpoints []string) string {
	if endpoints == nil || len(endpoints) == 0 {
		return ""
	}
	n := atomic.AddUint64(&r.cnt, 1)
	index := n % uint64(len(endpoints))
	return endpoints[index]
}
