// 版权所有 2021 github.com/gagliardetto
// 本文件已被github.com/gagliardetto修改
//
// 版权所有 2020 dfuse Platform Inc.
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
	"github.com/streamingfast/logging"
	"go.uber.org/zap"
)

var zlog *zap.Logger
var traceEnabled = logging.IsTraceEnabled("gosolana", "github.com/go-enols/gosolana/ws")

func init() {
	logging.Register("github.com/go-enols/gosolana/ws", &zlog)
}
