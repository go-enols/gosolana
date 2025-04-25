package main

import (
	"context"

	"github.com/go-enols/go-log"

	"github.com/gagliardetto/solana-go/rpc"
	"github.com/go-enols/gosolana"
)

var (
	Proxy               = "http://127.0.0.1:7890"
	NetWork rpc.Cluster = rpc.MainNetBeta
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	opt := gosolana.Option{
		WsUrl:   NetWork.WS,
		RpcUrl:  NetWork.RPC,
		Proxy:   Proxy,
		WsProxy: Proxy,
	}
	wallet1, err := gosolana.NewWallet(ctx, opt)
	if err != nil {
		log.Fatal("创建钱包失败 ", err)
	}
	opt.HTTPClient = wallet1.HTTPClient // 可以复用现有的http client，也可以在自定义http client然后复用
	wallet2, err := gosolana.NewWallet(ctx, opt)
	if err != nil {
		log.Fatal("创建钱包失败 ", err)
	}

	client1 := wallet1.GetClient()
	client2 := wallet2.GetClient()

	out, err := client1.GetBalance(ctx, wallet1.PublicKey(), rpc.CommitmentProcessed)
	if err != nil {
		log.Errorf("查询余额失败 | %s", err)
		return
	}
	log.Debug(out)

	out2, err := client2.GetBalance(ctx, wallet1.PublicKey(), rpc.CommitmentProcessed)
	if err != nil {
		log.Errorf("查询余额失败 | %s", err)
		return
	}
	log.Debug(out2)
}
