package gosolana

import (
	"context"
	"encoding/binary"
	"errors"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

type TokenMetaOnChain struct {
	Key                  uint8
	UpdateAuthority      solana.PublicKey
	Mint                 solana.PublicKey
	Name                 string
	Symbol               string
	Uri                  string
	SellerFeeBasisPoints uint16
	PrimarySaleHappened  uint8
	IsMutable            uint8
	EditionNonce         uint8
	TokenStandard        uint8
}

// 解析Metaplex Metadata账户数据
func ParseMetaplexMetadata(data []byte) (*TokenMetaOnChain, error) {
	if len(data) < 1+32+32+4+32+4+10+4+200+2+1+1+1+1 {
		return nil, errors.New("metadata data too short")
	}
	meta := &TokenMetaOnChain{}
	offset := 0

	meta.Key = data[offset]
	offset += 1

	copy(meta.UpdateAuthority[:], data[offset:offset+32])
	offset += 32

	copy(meta.Mint[:], data[offset:offset+32])
	offset += 32

	// name
	nameLen := int(binary.LittleEndian.Uint32(data[offset : offset+4]))
	offset += 4
	meta.Name = string(data[offset : offset+nameLen])
	offset += nameLen

	// symbol
	symbolLen := int(binary.LittleEndian.Uint32(data[offset : offset+4]))
	offset += 4
	meta.Symbol = string(data[offset : offset+symbolLen])
	offset += symbolLen

	// uri
	uriLen := int(binary.LittleEndian.Uint32(data[offset : offset+4]))
	offset += 4
	meta.Uri = string(data[offset : offset+uriLen])
	offset += uriLen

	meta.SellerFeeBasisPoints = binary.LittleEndian.Uint16(data[offset : offset+2])
	offset += 2

	meta.PrimarySaleHappened = data[offset]
	offset += 1

	meta.IsMutable = data[offset]
	offset += 1

	meta.EditionNonce = data[offset]
	offset += 1

	meta.TokenStandard = data[offset]
	// offset += 1

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
