package gosolana

import (
	"context"

	"log"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/go-enols/gosolana/ws"
)

type Wallet struct {
	rpc        *rpc.Client
	wsRpc      *ws.Client
	Address    string
	Base58Pkey string // base58格式的私钥
	HashPkey   string // hash格式的私钥

	*solana.Wallet
}

func NewWallet(ctx context.Context, option ...Option) (*Wallet, error) {
	op := NewDefaultOption(ctx, option...)
	var (
		wall *solana.Wallet
		err  error
	)
	if op.Pkey != "" {
		wall, err = solana.WalletFromPrivateKeyBase58(op.Pkey)
		if err != nil {
			return nil, err
		}
	}

	log.Printf("[DEBUG] 成功创建Solana钱包 | %s | %s", wall.PublicKey(), wall.PrivateKey)
	return &Wallet{
		rpc:        op.RpcClient,
		Address:    wall.PublicKey().String(),
		Base58Pkey: wall.PrivateKey.String(),
		HashPkey:   hexutil.Encode(wall.PrivateKey),
		Wallet:     wall,
		wsRpc:      op.WsClient,
	}, nil
}

func (w *Wallet) GetClient() *rpc.Client {
	return w.rpc
}

func (w *Wallet) GetWsClient() *ws.Client {
	return w.wsRpc
}

func (w *Wallet) SendTransaction(ctx context.Context, instruction []solana.Instruction) (bool, error) {
	recentBlockHash, err := w.rpc.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		log.Printf("[ERROR] 获取Hash失败 | %s", err)
		return false, err
	}
	// 构造交易
	tx, err := solana.NewTransaction(
		instruction,
		recentBlockHash.Value.Blockhash,
		solana.TransactionPayer(w.PublicKey()),
	)
	if err != nil {
		log.Printf("[ERROR] 构建交易失败 | %s", err)
		return false, err
	}

	// 签名交易
	out, err := tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if key.Equals(w.PublicKey()) {
			return &w.Wallet.PrivateKey
		}
		return nil
	})
	if err != nil {
		log.Printf("[ERROR] 签名交易失败 | %s", err)
		return false, err
	}
	log.Printf("[DEBUG] 签名交易输出 | %v", out)
	// 7. 发送交易
	sig, err := w.rpc.SendTransactionWithOpts(
		context.Background(),
		tx,
		rpc.TransactionOpts{
			SkipPreflight:       false,
			PreflightCommitment: rpc.CommitmentFinalized,
		},
	)
	if err != nil {
		log.Printf("[ERROR] 发送交易失败 | %s", err)
		return false, err
	}
	log.Printf("[INFO] Transaction Signature: %s", sig)
	log.Printf("[INFO] 交易详情 | %v", tx) // 打印交易详情
	result, err := w.GetTransaction(ctx, sig)
	if err != nil {
		log.Printf("[ERROR] 获取交易状态失败 | %s", err)
		return false, err
	}
	return result, nil
}

// GetTransaction 获取交易状态直到成功为止
//
// ctx: 上下文对象，方便后续设置超时等信息
//
// sign: 交易广播的sign
func (w *Wallet) GetTransaction(ctx context.Context, sign solana.Signature, option ...rpc.CommitmentType) (bool, error) {
	// 设置默认的确认等级 默认使用最快，也就是被单个服务器确认，但是还没有大量的服务器确认，即交易完成
	var commitment = rpc.CommitmentProcessed
	if len(option) > 0 {
		commitment = option[0]
	}
	// 等待交易确认
	sub, err := w.wsRpc.SignatureSubscribe(
		sign,
		commitment,
	)
	if err != nil {
		log.Printf("[ERROR] Failed to subscribe to signature: %v", err)
		return false, err
	}
	defer sub.Unsubscribe()

	for {
		got, err := sub.Recv(ctx)
		if err != nil {
			log.Printf("[ERROR] Error receiving signature status: %v", err)
			return false, err
		}
		if got.Value.Err != nil {
			log.Printf("[ERROR] Transaction failed: %v", got.Value.Err)
			return false, err
		} else {
			log.Printf("[INFO] Transaction confirmed | %s", sign.String())
			return true, nil
		}
	}
}
