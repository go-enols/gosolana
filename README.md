# gosolana - Solana WebSocket 订阅客户端

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Go 语言实现的 Solana 区块链 WebSocket 订阅客户端，基于 Solana 官方 API 规范开发，完全兼容官方接口规范，同时提供更高效的实时数据订阅功能和扩展特性。

## 功能特性

- 支持 Solana 主网和测试网的 WebSocket 连接
- 提供多种订阅类型：
  - 账户订阅
  - 区块订阅
  - 日志订阅
  - 程序订阅
  - 签名订阅
  - 投票订阅
- 自动重连机制
- 心跳检测保持连接活跃
- 线程安全的订阅管理

## 安装

```bash
go get github.com/go-enols/gosolana
```

## 快速开始

[订阅账户](./examples/subscribe_account/main.go)
[查询账户余额](./examples/get_balance/main.go)
[复用客户端](./examples/content/main.go)

## 实验性功能

[GetTokenAccount](./wallet.go#L172) 是由[helius](https://www.helius.dev/)提供的 Api 方法，如果你的 api 没有此功能你不应该调用他

## 贡献

欢迎提交 Pull Request 或 Issue 报告问题。

## 许可证

Apache 2.0 许可证
