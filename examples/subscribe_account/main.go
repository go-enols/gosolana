package main

import (
	"context"

	"github.com/davecgh/go-spew/spew"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/go-enols/gosolana"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// 连接到Solana主网
	wallet, err := gosolana.NewWallet(ctx, gosolana.Option{
		RpcUrl:  rpc.DevNet_RPC,          // 设置rpc连接（可选）
		WsUrl:   rpc.DevNet_WS,           // 设置ws连接（可选）
		Proxy:   "http://127.0.0.1:7890", // 设置http代理（可选）
		WsProxy: "http://127.0.0.1:7890", // 设置ws代理（可选）
	})
	if err != nil {
		panic(err)
	}

	client := wallet.GetWsClient()

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
