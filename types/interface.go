package types

import "context"

type Transaction interface {
	Commit(ctx context.Context) error
	Receipt() string
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