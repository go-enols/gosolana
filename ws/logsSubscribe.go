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

type LogResult struct {
	Context struct {
		Slot uint64
	} `json:"context"`
	Value struct {
		// The transaction signature.
		Signature solana.Signature `json:"signature"`
		// Error if transaction failed, null if transaction succeeded.
		Err interface{} `json:"err"`
		// Array of log messages the transaction instructions output
		// during execution, null if simulation failed before the transaction
		// was able to execute (for example due to an invalid blockhash
		// or signature verification failure)
		Logs []string `json:"logs"`
	} `json:"value"`
}

type LogsSubscribeFilterType string

const (
	// Subscribe to all transactions except for simple vote transactions.
	LogsSubscribeFilterAll LogsSubscribeFilterType = "all"
	// Subscribe to all transactions including simple vote transactions.
	LogsSubscribeFilterAllWithVotes LogsSubscribeFilterType = "allWithVotes"
)

// LogsSubscribe 订阅交易日志
func (cl *Client) LogsSubscribe(
	// Filter criteria for the logs to receive results by account type.
	filter LogsSubscribeFilterType,
	commitment rpc.CommitmentType, // (optional)
) (*LogSubscription, error) {
	return cl.logsSubscribe(
		filter,
		commitment,
	)
}

// LogsSubscribeMentions 订阅所有提及指定公钥的交易
func (cl *Client) LogsSubscribeMentions(
	// Subscribe to all transactions that mention the provided Pubkey.
	mentions solana.PublicKey,
	// (optional)
	commitment rpc.CommitmentType,
) (*LogSubscription, error) {
	return cl.logsSubscribe(
		rpc.M{
			"mentions": []string{mentions.String()},
		},
		commitment,
	)
}

// logsSubscribe 订阅交易日志(内部实现)
func (cl *Client) logsSubscribe(
	filter interface{},
	commitment rpc.CommitmentType,
) (*LogSubscription, error) {

	params := []interface{}{filter}
	conf := map[string]interface{}{}
	if commitment != "" {
		conf["commitment"] = commitment
	}

	genSub, err := cl.subscribe(
		params,
		conf,
		"logsSubscribe",
		"logsUnsubscribe",
		func(msg []byte) (interface{}, error) {
			var res LogResult
			err := decodeResponseFromMessage(msg, &res)
			return &res, err
		},
	)
	if err != nil {
		return nil, err
	}
	return &LogSubscription{
		sub: genSub,
	}, nil
}

type LogSubscription struct {
	sub *Subscription
}

func (sw *LogSubscription) Recv(ctx context.Context) (*LogResult, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case d, ok := <-sw.sub.stream:
		if !ok {
			return nil, ErrSubscriptionClosed
		}
		return d.(*LogResult), nil
	case err := <-sw.sub.err:
		return nil, err
	}
}

func (sw *LogSubscription) Err() <-chan error {
	return sw.sub.err
}

func (sw *LogSubscription) Response() <-chan *LogResult {
	typedChan := make(chan *LogResult, 1)
	go func(ch chan *LogResult) {
		// TODO: will this subscription yield more than one result?
		d, ok := <-sw.sub.stream
		if !ok {
			return
		}
		ch <- d.(*LogResult)
	}(typedChan)
	return typedChan
}

func (sw *LogSubscription) Unsubscribe() {
	sw.sub.Unsubscribe()
}
