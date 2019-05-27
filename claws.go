package claws

import (
	"context"
	"github.com/fadeAce/claws/types"
	"math/big"
)

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
	NewAddr() types.Bundle
	// return new addr
	BuildBundle(prv, pub, addr string) types.Bundle
	// return initiated TXN
	BuildTxn(hash string) types.TXN
	// deposit to certain place
	Withdraw(addr types.Bundle) *types.TxnInfo
	// seek for tx , keep track it
	Seek(txn types.TXN) bool

	// balance
	Balance(bundle types.Bundle) (string, error)

	// txs in block
	UnfoldTxs(ctx context.Context, num *big.Int) ([]types.TXN, error)

	// notify head is a blocked invoke
	NotifyHead(ctx context.Context, f func(num *big.Int)) error

	// Info used for info scan gap limit
	Info() *types.Info

}
