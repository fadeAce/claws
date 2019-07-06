package line

import (
	"context"
	types2 "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fadeAce/claws/types"
	"math/big"
	"sync"

	"github.com/fadeAce/claws/addr"
)

type erc20Wallet struct {
	*sync.RWMutex
	once sync.Once
	conf *types.Claws
	ctx  context.Context
	conn *Erc20Client
}

func NewERC20Wallet(conf *types.Claws, ctx context.Context, conn *Erc20Client, lock *sync.RWMutex) *erc20Wallet {
	res := &erc20Wallet{
		conf:    conf,
		ctx:     ctx,
		conn:    conn,
		RWMutex: lock,
	}
	return res
}

type Erc20Client struct {
	Conn   *ethclient.Client
	Closed bool
	header chan *types2.Header
	url    string
	ctx    context.Context
}

func NewERC20Client(ctx context.Context, url string) (client *Erc20Client, err error) {
	obj := &Erc20Client{
		header: make(chan *types2.Header),
		url:    url,
		ctx:    ctx,
		Closed: false,
	}

	c, err1 := ethclient.Dial(url)
	if err1 != nil {
		return nil, err1
	}

	obj.Conn = c
	_, err2 := c.SubscribeNewHead(ctx, obj.header)
	return obj, err2
}

func (c *Erc20Client) Reconnect() error {
	client, err1 := ethclient.Dial(c.url)
	if err1 != nil {
		return err1
	}

	c.Conn = client
	_, err2 := client.SubscribeNewHead(c.ctx, c.header)
	return err2
}

func (e *erc20Wallet) Send(ctx context.Context, from, to types.Bundle, amount string, option *types.Option) (tx types.Transaction, err error) {
	panic("implement me")
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
	//hash := txn.HexStr()
	//hs := common.HexToHash(hash)
	//reciept, err := e.conn.TransactionReceipt(e.ctx, hs)
	//if err != nil {
	//	log.Error("error", err)
	//	return false
	//}
	//if reciept.Status == types2.ReceiptStatusSuccessful {
	//	return true
	//}
	return false
}

// seek for tx , keep track it
func (e *erc20Wallet) Balance(bundle types.Bundle) (string, error) {
	//add := bundle.AddressStr()
	//address := common.HexToAddress(add)
	//balance, err := e.conn.BalanceAt(e.ctx, address, nil)
	//return balance.String(), err
	return "", nil
}

// seek for tx , keep track it
func (e *erc20Wallet) Type() string {
	return types.COIN_ERC20
}

func (e *erc20Wallet) UnfoldTxs(ctx context.Context, num *big.Int) (res []types.TXN, err error) {
	//b, err := e.conn.BlockByNumber(ctx, num)
	//if err != nil {
	//	return nil, err
	//}
	//txs := b.Transactions()
	//for _, v := range txs {
	//	f := v.GetFromUnsafe().String()
	//	txn := &erc20TXN{
	//		from:   f,
	//		Hash:   v.Hash().String(),
	//		fee:    new(big.Int).Mul(new(big.Int).SetUint64(v.Gas()), v.GasPrice()),
	//		amount: v.Value(),
	//	}
	//	if v.To() != nil {
	//		txn.to = v.To().String()
	//
	//		res = append(res, txn)
	//	} else {
	//		txn.to = ""
	//		res = append(res, txn)
	//	}
	//}
	//
	return
}

func (e *erc20Wallet) NotifyHead(ctx context.Context, fn func(num *big.Int)) (err error) {
	e.RLock()
	defer e.RUnlock()

	for {
		h := <-e.conn.header
		fn(h.Number)
	}

	return nil
}

func (e *erc20Wallet) Info() (info *types.Info) {
	return &types.Info{}
}
