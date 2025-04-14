// 版权所有 2022 github.com/gagliardetto
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
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

type BlockResult struct {
	Context struct {
		Slot uint64
	} `json:"context"`
	Value struct {
		Slot  uint64              `json:"slot"`
		Err   interface{}         `json:"err,omitempty"`
		Block *rpc.GetBlockResult `json:"block,omitempty"`
	} `json:"value"`
}

type BlockSubscribeFilter interface {
	isBlockSubscribeFilter()
}

var _ BlockSubscribeFilter = BlockSubscribeFilterAll("")

type BlockSubscribeFilterAll string

func (_ BlockSubscribeFilterAll) isBlockSubscribeFilter() {}

type BlockSubscribeFilterMentionsAccountOrProgram struct {
	Pubkey solana.PublicKey `json:"pubkey"`
}

func (_ BlockSubscribeFilterMentionsAccountOrProgram) isBlockSubscribeFilter() {}

func NewBlockSubscribeFilterAll() BlockSubscribeFilter {
	return BlockSubscribeFilterAll("")
}

func NewBlockSubscribeFilterMentionsAccountOrProgram(pubkey solana.PublicKey) *BlockSubscribeFilterMentionsAccountOrProgram {
	return &BlockSubscribeFilterMentionsAccountOrProgram{
		Pubkey: pubkey,
	}
}

type BlockSubscribeOpts struct {
	Commitment rpc.CommitmentType
	Encoding   solana.EncodingType `json:"encoding,omitempty"`

	// Level of transaction detail to return.
	TransactionDetails rpc.TransactionDetailsType

	// Whether to populate the rewards array. If parameter not provided, the default includes rewards.
	Rewards *bool

	// Max transaction version to return in responses.
	// If the requested block contains a transaction with a higher version, an error will be returned.
	MaxSupportedTransactionVersion *uint64
}

// 注意：不稳定功能，默认禁用
//
// 订阅新块确认或最终化的通知
//
// **此订阅功能不稳定，仅在验证器启动时添加
// `--rpc-pubsub-enable-block-subscription`标志时可用。
// 此订阅的格式在未来可能会发生变化**
func (cl *Client) BlockSubscribe(
	filter BlockSubscribeFilter,
	opts *BlockSubscribeOpts,
) (*BlockSubscription, error) {
	var params []interface{}
	if filter != nil {
		switch v := filter.(type) {
		case BlockSubscribeFilterAll:
			params = append(params, "all")
		case *BlockSubscribeFilterMentionsAccountOrProgram:
			params = append(params, rpc.M{"mentionsAccountOrProgram": v.Pubkey})
		}
	}
	if opts != nil {
		obj := make(rpc.M)
		if opts.Commitment != "" {
			obj["commitment"] = opts.Commitment
		}
		if opts.Encoding != "" {
			if !solana.IsAnyOfEncodingType(
				opts.Encoding,
				// Valid encodings:
				// solana.EncodingJSON, // TODO
				solana.EncodingJSONParsed, // TODO
				solana.EncodingBase58,
				solana.EncodingBase64,
				solana.EncodingBase64Zstd,
			) {
				return nil, fmt.Errorf("provided encoding is not supported: %s", opts.Encoding)
			}
			obj["encoding"] = opts.Encoding
		}
		if opts.TransactionDetails != "" {
			obj["transactionDetails"] = opts.TransactionDetails
		}
		if opts.Rewards != nil {
			obj["rewards"] = opts.Rewards
		}
		if opts.MaxSupportedTransactionVersion != nil {
			obj["maxSupportedTransactionVersion"] = *opts.MaxSupportedTransactionVersion
		}
		if len(obj) > 0 {
			params = append(params, obj)
		}
	}
	genSub, err := cl.subscribe(
		params,
		nil,
		"blockSubscribe",
		"blockUnsubscribe",
		func(msg []byte) (interface{}, error) {
			var res BlockResult
			err := decodeResponseFromMessage(msg, &res)
			return &res, err
		},
	)
	if err != nil {
		return nil, err
	}
	return &BlockSubscription{
		sub: genSub,
	}, nil
}

type BlockSubscription struct {
	sub *Subscription
}

func (sw *BlockSubscription) Recv(ctx context.Context) (*BlockResult, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case d, ok := <-sw.sub.stream:
		if !ok {
			return nil, ErrSubscriptionClosed
		}
		return d.(*BlockResult), nil
	case err := <-sw.sub.err:
		return nil, err
	}
}

func (sw *BlockSubscription) Err() <-chan error {
	return sw.sub.err
}

func (sw *BlockSubscription) Response() <-chan *BlockResult {
	typedChan := make(chan *BlockResult, 1)
	go func(ch chan *BlockResult) {
		// TODO: will this subscription yield more than one result?
		d, ok := <-sw.sub.stream
		if !ok {
			return
		}
		ch <- d.(*BlockResult)
	}(typedChan)
	return typedChan
}

func (sw *BlockSubscription) Unsubscribe() {
	sw.sub.Unsubscribe()
}
