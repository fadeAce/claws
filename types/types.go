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
	// chain type name
	Chain string
}

type Claws struct {
	Ctx     context.Context
	Version string `yaml:"version"`
	// The configuration represent coin.
	Coins []Coins `yaml:"coins"`

	// chain configuration of ethereum
	Eth *EthConf
}

type Coins struct {
	// Type of coin matched DB.
	CoinType string `yaml:"type"`
	// RPC location is configured to wallet builder
	// like 127.0.0.1:8545
	Url string `yml:"url"`

	ContractAddr string `yml:"contract_addr"`
}

type TxnInfo struct {
	Err    error
	TxHash string
	// "btc" "eth" "eth-c" from const
	TxType string
}

type Option struct {
	Nonce  uint64
	Secret string
	// fee here force decide the fee if it could be applied
	FeeConfig string
}

type EthConf struct {
	// name used for coins mapping
	Name string `yaml:"name"`
	// url used for setting client in RPC
	Url string `yaml:"url"`
}
