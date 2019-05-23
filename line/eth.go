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

type ethWallet struct {
	once sync.Once
	conf *types.Claws
	ctx  context.Context
	conn *ethclient.Client
}

type ethBundle struct {
	pub, prv, add string
}

type ethTXN struct {
	from   string
	to     string
	Hash   string
	amount *big.Int
	fee    *big.Int
}

func (etxn *ethTXN) HexStr() string {
	return etxn.Hash
}

func (etxn *ethTXN) SetStr(res string) {
	etxn.Hash = res
}

func (etxn *ethTXN) FromStr() string {
	return etxn.from
}

func (etxn *ethTXN) AmountStr() string {
	return etxn.amount.String()
}
func (etxn *ethTXN) FeeStr() string {
	return etxn.fee.String()
}

func (etxn *ethTXN) ToStr() string {
	return etxn.to
}

func (eb *ethBundle) InitAddr(pub, prv string) error {
	return nil
}

// return public key represented by hex encoding
func (eb *ethBundle) HexPubStr() string {
	return ""
}

// return private key represented by hex encoding
func (eb *ethBundle) HexPrvStr() string {
	return eb.prv
}

// return address in string
func (eb *ethBundle) AddressStr() string {
	return eb.add
}

func Type() string {
	return ""
}

// init wallet itself
func (e *ethWallet) InitWallet() {
}

// return new addr
func (e *ethWallet) NewAddr() types.Bundle {
	ad, err := addr.GenerateAddr(e.Type())
	if err != nil {
		return nil
	}
	res := &ethBundle{
		pub: "",
		prv: ad.HexPrvStr(),
		add: ad.HexAddrStr(),
	}
	return res
}

// return new addr
func (e *ethWallet) BuildBundle(
	prv, pub, addr string,
) types.Bundle {
	res := &ethBundle{
		pub: pub,
		prv: prv,
		add: addr,
	}
	return res
}

// return new addr
func (e *ethWallet) BuildTxn(
	tx string,
) types.TXN {
	return &ethTXN{
		Hash: tx,
	}
}

// withdraw to certain place
func (e *ethWallet) Withdraw(addr types.Bundle) *types.TxnInfo {
	return nil
}

// seek for tx , keep track it
func (e *ethWallet) Seek(txn types.TXN) bool {
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
func (e *ethWallet) Balance(bundle types.Bundle) (string, error) {
	add := bundle.AddressStr()
	address := common.HexToAddress(add)
	balance, err := e.conn.BalanceAt(e.ctx, address, nil)
	return balance.String(), err
}

// seek for tx , keep track it
func (e *ethWallet) Type() string {
	return types.COIN_ETH
}

func NewEthWallet(conf *types.Claws, ctx context.Context, conn *ethclient.Client) ethWallet {
	res := ethWallet{
		conf: conf,
		ctx:  ctx,
		conn: conn,
	}
	return res
}

func (e *ethWallet) UnfoldTxs(ctx context.Context, num *big.Int) (res []types.TXN, err error) {
	b, err := e.conn.BlockByNumber(ctx, num)
	if err != nil {
		return nil, err
	}
	txs := b.Transactions()
	for _, v := range txs {
		f := v.GetFromUnsafe().String()
		txn := &ethTXN{
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

func (e *ethWallet) NotifyHead(ctx context.Context, f func(num *big.Int)) (err error) {
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
