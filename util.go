package gosolana

import (
	"context"
	"encoding/binary"
	"errors"
	"github.com/gagliardetto/solana-go"
)

// TransferInstruction 示例：转账指令扩展
type TransferInstruction struct {
	BaseInstruction
	From   solana.PublicKey
	To     solana.PublicKey
	Amount uint64
}

func NewTransferInstruction(from, to solana.PublicKey, amount uint64) *TransferInstruction {
	return &TransferInstruction{
		BaseInstruction: BaseInstruction{
			programID: solana.SystemProgramID,
			accounts: []*solana.AccountMeta{
				{PublicKey: from, IsSigner: true, IsWritable: true},
				{PublicKey: to, IsSigner: false, IsWritable: true},
			},
			dataCoder: &TransferData{
				InstructionType: 2,
				Amount:          amount,
			},
		},
	}
}

// TransferData 转账数据编码器
type TransferData struct {
	InstructionType uint32 // 指令类型标识
	Amount          uint64
}

func (td *TransferData) Decode(data []byte) error {
	// 基础长度校验（4字节指令类型+8字节金额）
	if len(data) < 12 {
		return errors.New("invalid data length")
	}
	// 分离指令类型和金额数据
	td.InstructionType = binary.LittleEndian.Uint32(data[0:4])
	td.Amount = binary.LittleEndian.Uint64(data[4:12])
	return nil
}

func (td *TransferData) Encode() ([]byte, error) {
	buf := make([]byte, 12)
	binary.LittleEndian.PutUint32(buf[0:4], td.InstructionType)
	binary.LittleEndian.PutUint64(buf[4:12], td.Amount)
	return buf, nil
}

// SendSol 向某个钱包发送sol 如果需要大量的转移应该重新构造交易指令，以事务的方式转移
func (w *Wallet) SendSol(ctx context.Context, to string, amount float64) {
	fromAccount := solana.MustPublicKeyFromBase58(w.PublicKey().String())
	toAccount := solana.MustPublicKeyFromBase58(to)
	toAmount := uint64(amount * float64(solana.LAMPORTS_PER_SOL))

	if _, err := w.SendTransaction(ctx, []solana.Instruction{
		// 构造交易指令
		NewTransferInstruction(fromAccount, toAccount, toAmount),
	}); err != nil {
		return
	}
}
