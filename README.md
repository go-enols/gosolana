# gosolana - Solana WebSocket订阅客户端

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Go语言实现的Solana区块链WebSocket订阅客户端，基于Solana官方API规范开发，完全兼容官方接口规范，同时提供更高效的实时数据订阅功能和扩展特性。

## 功能特性

- 支持Solana主网和测试网的WebSocket连接
- 提供多种订阅类型：
  - 账户订阅
  - 区块订阅
  - 日志订阅
  - 程序订阅
  - 签名订阅
  - 投票订阅
- 自动重连机制
- 心跳检测保持连接活跃
- 线程安全的订阅管理

## 安装

```bash
go get github.com/go-enols/gosolana
```

## 快速开始

### 账户订阅示例

```go
package main

import (
	"context"

	"github.com/davecgh/go-spew/spew"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
)

func main() {
	// 连接到Solana主网
	client, err := ws.Connect(context.Background(), gosolana.NewWallet(ctx, gosolana.Option{
		RpcUrl:  rpc.DevNet_RPC, // 设置rpc连接（可选）
		WsUrl:   rpc.DevNet_WS, // 设置ws连接（可选）
		Proxy:   "http://127.0.0.1:7890", // 设置http代理（可选）
		WsProxy: "http://127.0.0.1:7890", // 设置ws代理（可选）
	})
	if err != nil {
		panic(err)
	}
	defer client.Close()
	
	// 订阅Serum程序账户
	program := solana.MustPublicKeyFromBase58("9xQeWvG816bUx9EPjHmaT23yvVM2ZWbrrpZb9PusVFin")
	
	sub, err := client.AccountSubscribe(program, "")
	if err != nil {
		panic(err)
	}
	defer sub.Unsubscribe()

	// 接收账户变更通知
	for {
		got, err := sub.Recv(context.Background())
		if err != nil {
			panic(err)
		}
		spew.Dump(got)
	}
}
```

### 其他订阅类型

项目还提供以下订阅功能，使用方法类似：

- `BlockSubscribe` - 订阅新区块
- `LogsSubscribe` - 订阅日志
- `ProgramSubscribe` - 订阅程序账户变更
- `SignatureSubscribe` - 订阅交易签名状态
- `VoteSubscribe` - 订阅投票交易

## API参考

### 连接管理

- `Connect(ctx, endpoint)` - 创建WebSocket连接
- `ConnectWithOptions(ctx, endpoint, options)` - 带选项创建连接
- `Close()` - 关闭连接

### 订阅管理

- `AccountSubscribe(account, commitment)` - 订阅账户变更
- `BlockSubscribe(commitment)` - 订阅新区块
- `LogsSubscribe(filter, commitment)` - 订阅日志
- `ProgramSubscribe(program, commitment)` - 订阅程序账户变更
- `SignatureSubscribe(signature, commitment)` - 订阅交易签名状态
- `VoteSubscribe()` - 订阅投票交易

## 贡献

欢迎提交Pull Request或Issue报告问题。

## 许可证

Apache 2.0 许可证
