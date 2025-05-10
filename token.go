package gosolana

import (
	"context"
	"errors"
	"fmt"
	"strings"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"

	token_metadata "github.com/go-enols/metaplex-go/clients/token-metadata"
)

// 获取代币的元数据
func GetTokenMetaOnChain(ctx context.Context, client *rpc.Client, mint solana.PublicKey) (*token_metadata.Metadata, error) {
	metaPDA, err := GetMetadata(mint)
	if err != nil {
		return nil, fmt.Errorf("derive metadata pda failed: %w", err)
	}

	info, err := client.GetAccountInfo(ctx, metaPDA)
	if err != nil {
		return nil, fmt.Errorf("get metadata account failed: %w", err)
	}
	if info.Value == nil || info.Value.Data == nil {
		return nil, errors.New("metadata account not found")
	}

	res := new(token_metadata.Metadata)
	if err := bin.NewDecoderWithEncoding(info.GetBinary(), bin.EncodingBorsh).Decode(res); err != nil {
		return nil, err
	}

	//\x00
	res.Data.Name = strings.Trim(res.Data.Name, "\x00")
	res.Data.Symbol = strings.Trim(res.Data.Symbol, "\x00")

	return res, nil
}

// 获取代币的SPL派生钱包
//
// 例如：获取自己账户(wallet)在代币USDCT(mint)下的衍生钱包
//
// 用例 接收、交换、铸造代币时都需要
func GetTokenAccount(wallet solana.PublicKey, mint solana.PublicKey) (solana.PublicKey, error) {
	addr, _, err := solana.FindProgramAddress(
		[][]byte{
			wallet.Bytes(),
			solana.TokenProgramID.Bytes(),
			mint.Bytes(),
		},
		solana.SPLAssociatedTokenAccountProgramID,
	)
	return addr, err
}

// 获取metadata的账户
func GetMetadata(mint solana.PublicKey) (solana.PublicKey, error) {
	addr, _, err := solana.FindProgramAddress(
		[][]byte{
			[]byte("metadata"),
			token_metadata.ProgramID.Bytes(),
			mint.Bytes(),
		},
		token_metadata.ProgramID,
	)
	return addr, err
}
