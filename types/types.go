package types

import "context"

const (
	COIN_BTC = "btc"
	COIN_ETH = "eth"
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

type Marble struct {
	Ctx     context.Context
	Version string `yaml:"version"`
	// The configuration represent coin.
	Coins []struct {
		// The estimate fee for this type of coin. This is usually represented by it's own coin convention.
		Fee string `yaml:"fee"`

		// Represent coin to be charged.
		FeeCoin string `yaml:"fee_coin"`

		// Full name of this type of coin.
		FullName string `yaml:"full_name"`

		// Abbreviated name of this type of coin.
		ShortName string `yaml:"short_name"`

		// Type of coin matched DB.
		CoinType string `yaml:"type"`

		// How precise it would be for decimal balance.
		Decimal int `yaml:"decimal"`

		// minimum withdraw limitation AKA : gap limit
		GapLimit float64

		// RPC location is configured to wallet builder
		// like 127.0.0.1:8545
		Url string `yml:"url"`

		// interval marks the interval scan time for each chain
		Interval int `yml:"interval"`
	} `yaml:"coins"`
}

type TxnInfo struct {
	Err    error
	TxHash string
	// "btc" "eth" "eth-c" from const
	TxType string
}
