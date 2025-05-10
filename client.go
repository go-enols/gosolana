package gosolana

import (
	"context"

	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/jsonrpc"
	"github.com/go-enols/go-log"
)

type ClientOption func(ctx context.Context, client *client)

// 从原先的Wallet中改造过来,进行负载均衡
type client struct {
	ctx context.Context // 父级上下文

	rpc []*rpc.Client
}

// TODO:还未完成先搁置，后面需要转移到这里，完成负载均衡
func NewRPCClient(ctx context.Context, netWork rpc.Cluster, opt ...ClientOption) *client {
	c := new(client)
	c.ctx = ctx

	for _, fn := range opt {
		fn(ctx, c)
	}
	if len(c.rpc) == 0 { // 如果没有rpc就填入默认节点
		defaultRPC(netWork) // 创建默认的RPC节点,如果用户自定义了就直接覆盖
	}

	return c
}

// 默认的RPC节点
func defaultRPC(netWork rpc.Cluster) ClientOption {
	return func(ctx context.Context, client *client) {
		client.rpc = []*rpc.Client{
			rpc.NewWithCustomRPCClient(
				rpc.NewWithRateLimit(netWork.RPC, 1), // 设置请求限制，每秒1条
			),
		}
	}
}

// 设置RPC代理
//
//	endpoint 节点地址
//	rps 速度限制每秒请求多少次（可选）
func WithRPCProxy(endpoint, proxy string, rps ...int) ClientOption {
	rp := 0
	if len(rps) > 0 {
		rp = rps[0]
	}
	httpClient, err := NewProxyHttpClient(proxy)
	if err != nil {
		log.Error("你提供了一个无效的代理")
		return func(ctx context.Context, client *client) {

		}
	}

	return func(ctx context.Context, client *client) {
		if rp != 0 { // 如果设置了速度限制
			client.rpc = append(client.rpc, rpc.NewWithCustomRPCClient(
				NewWithRateLimit(endpoint, rp, &jsonrpc.RPCClientOpts{
					HTTPClient: httpClient,
				}),
			))
		} else { // 不进行速度限制
			client.rpc = append(client.rpc,
				rpc.NewWithCustomRPCClient(
					jsonrpc.NewClientWithOpts(endpoint, &jsonrpc.RPCClientOpts{
						HTTPClient: httpClient,
					})),
			)
		}
	}
}

// 设置一个已有的rpc节点
func WithRPCClient(rpc *rpc.Client) ClientOption {
	return func(ctx context.Context, client *client) {
		client.rpc = append(client.rpc, rpc)
	}
}
