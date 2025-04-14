// 版权所有 2021 github.com/gagliardetto
//
// 根据Apache许可证2.0版授权
// 除非遵守许可证，否则不得使用此文件
// 您可以在以下网址获取许可证副本
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// 除非适用法律要求或书面同意，软件
// 根据许可证分发是基于"按原样"的基础，
// 没有任何明示或暗示的保证或条件
// 请参阅许可证了解特定语言的权限和限制

package ws

import (
	"context"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

type ProgramResult struct {
	Context struct {
		Slot uint64
	} `json:"context"`
	Value rpc.KeyedAccount `json:"value"`
}

// ProgramSubscribe 订阅程序账户变更通知
// 当程序拥有的账户lamports或数据发生变化时接收通知
func (cl *Client) ProgramSubscribe(
	programID solana.PublicKey,
	commitment rpc.CommitmentType,
) (*ProgramSubscription, error) {
	return cl.ProgramSubscribeWithOpts(
		programID,
		commitment,
		"",
		nil,
	)
}

// ProgramSubscribeWithOpts 带选项订阅程序账户变更通知
// 当程序拥有的账户lamports或数据发生变化时接收通知
func (cl *Client) ProgramSubscribeWithOpts(
	programID solana.PublicKey,
	commitment rpc.CommitmentType,
	encoding solana.EncodingType,
	filters []rpc.RPCFilter,
) (*ProgramSubscription, error) {

	params := []interface{}{programID.String()}
	conf := map[string]interface{}{
		"encoding": "base64",
	}
	if commitment != "" {
		conf["commitment"] = commitment
	}
	if encoding != "" {
		conf["encoding"] = encoding
	}
	if filters != nil && len(filters) > 0 {
		conf["filters"] = filters
	}

	genSub, err := cl.subscribe(
		params,
		conf,
		"programSubscribe",
		"programUnsubscribe",
		func(msg []byte) (interface{}, error) {
			var res ProgramResult
			err := decodeResponseFromMessage(msg, &res)
			return &res, err
		},
	)
	if err != nil {
		return nil, err
	}
	return &ProgramSubscription{
		sub: genSub,
	}, nil
}

type ProgramSubscription struct {
	sub *Subscription
}

func (sw *ProgramSubscription) Recv(ctx context.Context) (*ProgramResult, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case d, ok := <-sw.sub.stream:
		if !ok {
			return nil, ErrSubscriptionClosed
		}
		return d.(*ProgramResult), nil
	case err := <-sw.sub.err:
		return nil, err
	}
}

func (sw *ProgramSubscription) Err() <-chan error {
	return sw.sub.err
}

func (sw *ProgramSubscription) Response() <-chan *ProgramResult {
	typedChan := make(chan *ProgramResult, 1)
	go func(ch chan *ProgramResult) {
		// TODO: will this subscription yield more than one result?
		d, ok := <-sw.sub.stream
		if !ok {
			return
		}
		ch <- d.(*ProgramResult)
	}(typedChan)
	return typedChan
}

func (sw *ProgramSubscription) Unsubscribe() {
	sw.sub.Unsubscribe()
}
