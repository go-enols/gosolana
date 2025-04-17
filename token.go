package gosolana

import (
	"context"
	"encoding/binary"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

type TokenMetadata struct {
	MintAddress string `json:"mint"`
	Symbol      string `json:"symbol"`
	Name        string `json:"name"`
	Decimals    uint8  `json:"decimals"`
	TotalSupply uint64 `json:"total_supply"` // 新增字段
}

func GetTokenMetadata(client *rpc.Client, mintAddress solana.PublicKey) (TokenMetadata, error) {
	// 获取代币基本信息
	mintInfo, err := client.GetAccountInfo(context.TODO(), mintAddress)
	if err != nil {
		return TokenMetadata{}, err
	}

	var meta TokenMetadata
	meta.MintAddress = mintAddress.String()
	if mintInfo.Value != nil && mintInfo.Value.Data != nil {
		dec := mintInfo.Value.Data.GetBinary()
		if len(dec) >= 45 {
			// 解析decimals
			meta.Decimals = dec[44]
			// 解析total supply（前8字节，uint64，小端序）
			meta.TotalSupply = binary.LittleEndian.Uint64(dec[36:44])
		}
	}

	// 获取Metaplex扩展元数据
	metadataAccount, err := deriveMetadataPDA(mintAddress)
	if err != nil {
		return meta, err
	}
	metadata, err := client.GetAccountInfo(context.TODO(), metadataAccount)
	if err == nil && metadata.Value != nil && metadata.Value.Data != nil {
		name, symbol := parseMetaplexMetadata(metadata.Value.Data.GetBinary())
		meta.Name = name
		meta.Symbol = symbol
	}

	return meta, nil
}

// 查询多个 SPL Token 账户余额（返回 UiAmount 数组，顺序与输入一致）
func GetMultipleAccountsBalances(ctx context.Context, client *rpc.Client, accounts []solana.PublicKey) ([]*rpc.GetTokenAccountBalanceResult, error) {
	results := make([]*rpc.GetTokenAccountBalanceResult, len(accounts))
	for i, acc := range accounts {
		resp, err := client.GetTokenAccountBalance(ctx, acc, rpc.CommitmentFinalized)
		if err != nil {
			return nil, err
		}
		results[i] = resp
	}
	return results, nil
}
