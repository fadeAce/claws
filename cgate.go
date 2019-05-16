package claws

import "github.com/fadeAce/claws/types"

// gate have to access chain to get specific data by invoking RPC nodes

/*

			+------------------------------+
account --+ | wallet --|-- addr --|-- line | ---+ ETH/BTC
			+------------------------------+

*/

// wallet is general structure keep up with the chain in using line handler
type Wallet interface {
	// typ
	Type() string
	// init wallet itself
	InitWallet()
	// return new addr
	NewAddr() Bundle
	// return new addr
	BuildBundle(prv, pub, addr string) Bundle
	// return initiated TXN
	BuildTxn(hash string) TXN
	// deposit to certain place
	Withdraw(addr Bundle) *types.TxnInfo
	// seek for tx , keep track it
	Seek(txn TXN) bool
	// Info
	Info() *types.Info

	// balance
	Balance(bundle Bundle) (string, error)
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
}

type SeekInfo struct {
}

type WalletInfo struct {
}

type Require struct {
}
