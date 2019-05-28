package line

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	types2 "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/fadeAce/claws/types"
	"math/big"
	"sync"

	"github.com/fadeAce/claws/addr"
)

type erc20Wallet struct {
	once sync.Once
	conf *types.Claws
	ctx  context.Context
	conn *ethclient.Client
}

type erc20Bundle struct {
	pub, prv, add string
}

type erc20TXN struct {
	from   string
	to     string
	Hash   string
	amount *big.Int
	fee    *big.Int
}

func (etxn *erc20TXN) HexStr() string {
	return etxn.Hash
}

func (etxn *erc20TXN) SetStr(res string) {
	etxn.Hash = res
}

func (etxn *erc20TXN) FromStr() string {
	return etxn.from
}

func (etxn *erc20TXN) AmountStr() string {
	return etxn.amount.String()
}
func (etxn *erc20TXN) FeeStr() string {
	return etxn.fee.String()
}

func (etxn *erc20TXN) ToStr() string {
	return etxn.to
}

func (eb *erc20Bundle) InitAddr(pub, prv string) error {
	return nil
}

// return public key represented by hex encoding
func (eb *erc20Bundle) HexPubStr() string {
	return ""
}

// return private key represented by hex encoding
func (eb *erc20Bundle) HexPrvStr() string {
	return eb.prv
}

// return address in string
func (eb *erc20Bundle) AddressStr() string {
	return eb.add
}

// init wallet itself
func (e *erc20Wallet) InitWallet() {
}

// return new addr
func (e *erc20Wallet) NewAddr() types.Bundle {
	ad, err := addr.GenerateAddr(e.Type())
	if err != nil {
		return nil
	}
	res := &erc20Bundle{
		pub: "",
		prv: ad.HexPrvStr(),
		add: ad.HexAddrStr(),
	}
	return res
}

// return new addr
func (e *erc20Wallet) BuildBundle(
	prv, pub, addr string,
) types.Bundle {
	res := &erc20Bundle{
		pub: pub,
		prv: prv,
		add: addr,
	}
	return res
}

// return new addr
func (e *erc20Wallet) BuildTxn(
	tx string,
) types.TXN {
	return &erc20TXN{
		Hash: tx,
	}
}

// withdraw to certain place
func (e *erc20Wallet) Withdraw(addr types.Bundle) *types.TxnInfo {
	return nil
}

// seek for tx , keep track it
func (e *erc20Wallet) Seek(txn types.TXN) bool {
	// seek the txn hash of it
	hash := txn.HexStr()
	hs := common.HexToHash(hash)
	reciept, err := e.conn.TransactionReceipt(e.ctx, hs)
	if err != nil {
		log.Error("error", err)
		return false
	}
	if reciept.Status == types2.ReceiptStatusSuccessful {
		return true
	}
	return false
}

// seek for tx , keep track it
func (e *erc20Wallet) Balance(bundle types.Bundle) (string, error) {
	add := bundle.AddressStr()
	address := common.HexToAddress(add)
	balance, err := e.conn.BalanceAt(e.ctx, address, nil)
	return balance.String(), err
}

// seek for tx , keep track it
func (e *erc20Wallet) Type() string {
	return types.COIN_BTC
}

func Newerc20Wallet(conf *types.Claws, ctx context.Context, conn *ethclient.Client) erc20Wallet {
	res := erc20Wallet{
		conf: conf,
		ctx:  ctx,
		conn: conn,
	}
	return res
}

func (e *erc20Wallet) UnfoldTxs(ctx context.Context, num *big.Int) (res []types.TXN, err error) {
	b, err := e.conn.BlockByNumber(ctx, num)
	if err != nil {
		return nil, err
	}
	txs := b.Transactions()
	for _, v := range txs {
		f := v.GetFromUnsafe().String()
		txn := &erc20TXN{
			from:   f,
			Hash:   v.Hash().String(),
			fee:    new(big.Int).Mul(new(big.Int).SetUint64(v.Gas()), v.GasPrice()),
			amount: v.Value(),
		}
		if v.To() != nil {
			txn.to = v.To().String()

			res = append(res, txn)
		} else {
			txn.to = ""
			res = append(res, txn)
		}
	}

	return
}

func (e *erc20Wallet) NotifyHead(ctx context.Context, f func(num *big.Int)) (err error) {
	ch := make(chan *types2.Header)
	e.once.Do(func() {
		_, err = e.conn.SubscribeNewHead(ctx, ch)
		if err != nil {
			return
		}
		for {
			head := <-ch
			f(head.Number)
		}
	})
	return
}


func (e *erc20Wallet) Info() (info *types.Info) {
	return &types.Info{}
}
