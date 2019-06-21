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

// Addr marks for interface based addr
type Bundle interface {
	InitAddr(pub, prv string) error
	// return public key represented by hex encoding
	HexPubStr() string
	// return private key represented by hex encoding
	HexPrvStr() string
	// return address in string
	AddressStr() string
}

type TXN interface {
	HexStr() string
	SetStr(res string)
	FromStr() string
	ToStr() string
	FeeStr() string
	AmountStr() string
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
}
