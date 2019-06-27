package types

import (
	"context"
)

const (
	COIN_BTC   = "btc"
	COIN_ETH   = "eth"
	COIN_ERC20 = "erc20"
)

type Info struct {
	Name         string
	Fee          string
	DisplayShort string
	Display      string
	Decimal      int
	FeeCoin      string
	GapLimit     float64
}

type Claws struct {
	Ctx     context.Context
	Version string `yaml:"version"`
	// The configuration represent coin.
	Coins []Coins `yaml:"coins"`
}

type Coins struct {
	// Type of coin matched DB.
	CoinType string `yaml:"type"`
	// RPC location is configured to wallet builder
	// like 127.0.0.1:8545
	Url string `yml:"url"`
}

type TxnInfo struct {
	Err    error
	TxHash string
	// "btc" "eth" "eth-c" from const
	TxType string
}

type SeekInfo struct {
}

type WalletInfo struct {
}

type Require struct {
}

type Option struct {
	Nonce  uint64
	Secret string
	// fee here force decide the fee if it could be applied
	FeeConfig string
}
