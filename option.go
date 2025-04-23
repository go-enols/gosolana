package gosolana

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/go-enols/gosolana/ws"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/jsonrpc"
)

type Option struct {
	JsonRpcClient jsonrpc.RPCClient
	RpcClient     *rpc.Client
	WsClient      *ws.Client
	WsUrl         string
	RpcUrl        string
	Pkey          string
	Headers       map[string]string
	HTTPClient    *http.Client
	Proxy         string
	WsProxy       string
	TimeOut       time.Duration
}

// NewDefaultOption 构建一个新的配置项
func NewDefaultOption(ctx context.Context, option ...Option) Option {
	result := Option{
		Headers: make(map[string]string),
	}
	if len(option) > 0 {
		result = option[0]
	}
	if result.RpcUrl == "" {
		result.RpcUrl = rpc.DevNet_RPC
	}
	if result.WsUrl == "" {
		result.WsUrl = rpc.DevNet_WS
	}
	if len(result.Headers) == 0 {
		result.Headers = map[string]string{}
	}

	if result.HTTPClient == nil {
		if result.Proxy != "" {
			client, err := NewProxyHttpClient(result.Proxy)
			if err != nil {
				panic(err)
			}
			result.HTTPClient = client
		} else {
			result.HTTPClient = &http.Client{}
		}
	}
	if result.TimeOut == 0 {
		result.TimeOut = 5 * time.Second
	}
	// 如果用户没有设置请求超时，则默认请求5秒后超时
	result.HTTPClient.Timeout = result.TimeOut

	if result.RpcClient == nil {
		result.RpcClient = rpc.NewWithCustomRPCClient(jsonrpc.NewClientWithOpts(result.RpcUrl, &jsonrpc.RPCClientOpts{
			HTTPClient:    result.HTTPClient,
			CustomHeaders: result.Headers,
		}))
		result.JsonRpcClient = jsonrpc.NewClientWithOpts(result.RpcUrl, &jsonrpc.RPCClientOpts{
			HTTPClient: result.HTTPClient,
		})
	}
	if result.WsClient == nil {
		wsClient, err := ws.ConnectWithOptions(ctx, result.WsUrl, &ws.Options{
			Proxy: result.WsProxy,
		})
		if err != nil {
			panic(err)
		}
		result.WsClient = wsClient
	}

	if result.Pkey == "" {
		temp := solana.NewWallet()
		result.Pkey = temp.PrivateKey.String()
	}

	return result
}

// NewProxyHttpClient 创建一个支持代理的HTTP/HTTPS客户端
func NewProxyHttpClient(proxy string) (*http.Client, error) {
	proxyURL, err := url.Parse(proxy)
	if err != nil {
		return nil, fmt.Errorf("解析代理URL失败: %w", err)
	}

	// 克隆默认Transport以保留其他配置（如TLS、连接池）
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.Proxy = http.ProxyURL(proxyURL)

	return &http.Client{
		Transport: transport,
	}, nil
}

// NewDefaultRpcClient 创建一个适用于solana的JSONRPC 客户端他使用系统代理
func NewDefaultRpcClient(endpoint string, opts *jsonrpc.RPCClientOpts) *rpc.Client {
	return rpc.NewWithCustomRPCClient(jsonrpc.NewClientWithOpts(endpoint, opts))
}

// NewProxyRpcClient 创建一个适用于solana的JSONRPC 客户端他使用指定的代理
func NewProxyRpcClient(endpoint, proxy string) (*rpc.Client, error) {
	client, err := NewProxyHttpClient(proxy)
	if err != nil {
		return nil, err
	}
	return NewDefaultRpcClient(endpoint, &jsonrpc.RPCClientOpts{HTTPClient: client}), nil
}

// NewDefaultWsClient 创建一个默认配置的ws jsonrpc客户端 你可以指定他的配置
func NewDefaultWsClient(ctx context.Context, endpoint string, options *ws.Options) (*ws.Client, error) {
	return ws.ConnectWithOptions(ctx, endpoint, options)
}

// NewProxyWsClient 创建一个根据用户提供的代理地址的ws客户端
func NewProxyWsClient(ctx context.Context, endpoint, proxy string) (*ws.Client, error) {
	return NewDefaultWsClient(ctx, endpoint, &ws.Options{Proxy: proxy})
}
