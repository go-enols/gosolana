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
	wallet, err := gosolana.NewWallet(ctx, opt)
	if err != nil {
		log.Fatal("创建钱包失败 ", err)
	}

	client := wallet.GetClient()

	out, err := client.GetBalance(ctx, wallet.PublicKey(), rpc.CommitmentProcessed)
	if err != nil {
		log.Errorf("查询余额失败 | %s", err)
		return
	}
	log.Debug(out)
}
