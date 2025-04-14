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

type AccountResult struct {
	Context struct {
		Slot uint64
	} `json:"context"`
	Value struct {
		rpc.Account
	} `json:"value"`
}

// AccountSubscribe 订阅账户变更通知
// 当指定公钥账户的lamports或数据发生变化时接收通知
func (cl *Client) AccountSubscribe(
	account solana.PublicKey,
	commitment rpc.CommitmentType,
) (*AccountSubscription, error) {
	return cl.AccountSubscribeWithOpts(
		account,
		commitment,
		"",
	)
}

// AccountSubscribeWithOpts 带选项订阅账户变更通知
// 当指定公钥账户的lamports或数据发生变化时接收通知
func (cl *Client) AccountSubscribeWithOpts(
	account solana.PublicKey,
	commitment rpc.CommitmentType,
	encoding solana.EncodingType,
) (*AccountSubscription, error) {

	params := []interface{}{account.String()}
	conf := map[string]interface{}{
		"encoding": "base64",
	}
	if commitment != "" {
		conf["commitment"] = commitment
	}
	if encoding != "" {
		conf["encoding"] = encoding
	}

	genSub, err := cl.subscribe(
		params,
		conf,
		"accountSubscribe",
		"accountUnsubscribe",
		func(msg []byte) (interface{}, error) {
			var res AccountResult
			err := decodeResponseFromMessage(msg, &res)
			return &res, err
		},
	)
	if err != nil {
		return nil, err
	}
	return &AccountSubscription{
		sub: genSub,
	}, nil
}

type AccountSubscription struct {
	sub *Subscription
}

func (sw *AccountSubscription) Recv(ctx context.Context) (*AccountResult, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case d, ok := <-sw.sub.stream:
		if !ok {
			return nil, ErrSubscriptionClosed
		}
		return d.(*AccountResult), nil
	case err := <-sw.sub.err:
		return nil, err
	}
}

func (sw *AccountSubscription) Err() <-chan error {
	return sw.sub.err
}

func (sw *AccountSubscription) Response() <-chan *AccountResult {
	typedChan := make(chan *AccountResult, 1)
	go func(ch chan *AccountResult) {
		// TODO: will this subscription yield more than one result?
		d, ok := <-sw.sub.stream
		if !ok {
			return
		}
		ch <- d.(*AccountResult)
	}(typedChan)
	return typedChan
}

func (sw *AccountSubscription) Unsubscribe() {
	sw.sub.Unsubscribe()
}
