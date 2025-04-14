package gosolana

import (
	"github.com/gagliardetto/solana-go"
)

// Instruction 基础指令结构体
// 采用接口+实现类的设计模式，支持多类型指令扩展
type Instruction interface {
	ProgramID() solana.PublicKey    // 获取关联程序地址
	Accounts() []solana.AccountMeta // 获取账户元数据列表
	Data() ([]byte, error)          // 序列化指令数据
}

// BaseInstruction 基础指令实现
// 通过组合模式实现扩展性
type BaseInstruction struct {
	programID solana.PublicKey
	accounts  []*solana.AccountMeta
	data      []byte
	dataCoder DataCoder // 数据编码器接口
}

// ProgramID 实现Instruction接口
func (bi *BaseInstruction) ProgramID() solana.PublicKey {
	return bi.programID
}

// Accounts 实现Instruction接口
func (bi *BaseInstruction) Accounts() []*solana.AccountMeta {
	return bi.accounts
}

// Data 实现Instruction接口
func (bi *BaseInstruction) Data() ([]byte, error) {
	if bi.dataCoder != nil {
		return bi.dataCoder.Encode()
	}
	return bi.data, nil
}

// DataCoder 数据编码接口
// 支持自定义数据序列化逻辑
type DataCoder interface {
	Encode() ([]byte, error)
	Decode([]byte) error
}
