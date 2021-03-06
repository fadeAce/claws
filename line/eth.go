package line

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/fadeAce/claws/chain"
	"github.com/fadeAce/claws/types"
	"math/big"
	"strings"
	"sync"

	"github.com/fadeAce/claws/addr"
)

type ethWallet struct {
	*sync.RWMutex
	once sync.Once
	conf *types.Claws
	ctx  context.Context
	conn *chain.EthConn

	gasPrice *chain.EthGas

	updateCh chan interface{}
}

//type EthConn struct {
//	Conn *ethclient.Client
//}
//
//type EthGas struct {
//	GasPrice string
//}

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

type ethCommit struct {
	signTx *ethTypes.Transaction
	conn   *chain.EthConn
}

func (ec *ethCommit) Commit(ctx context.Context) error {
	err := ec.conn.Conn.SendTransaction(ctx, ec.signTx)
	if err != nil {
		return err
	}
	return nil
}

func (ec *ethCommit) Receipt() string {
	return ec.signTx.Hash().String()
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
	return strings.ToLower(eb.add)
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
		add: strings.ToLower(ad.HexAddrStr()),
	}
	return res
}

func NewConn(client *ethclient.Client) *chain.EthConn {
	return &chain.EthConn{client}
}

// return new addr
func (e *ethWallet) BuildBundle(
	prv, pub, addr string,
) types.Bundle {
	res := &ethBundle{
		pub: pub,
		prv: prv,
		add: strings.Trim(addr, " "),
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
	reciept, err := e.conn.Conn.TransactionReceipt(e.ctx, hs)
	if err != nil {
		log.Error("error", err)
		return false
	}
	if reciept.Status == ethTypes.ReceiptStatusSuccessful {
		return true
	}
	return false
}

// seek for tx , keep track it
func (e *ethWallet) Balance(bundle types.Bundle) (string, error) {
	add := bundle.AddressStr()
	address := common.HexToAddress(add)
	balance, err := e.conn.Conn.BalanceAt(e.ctx, address, nil)
	return balance.String(), err
}

// seek for tx , keep track it
func (e *ethWallet) Type() string {
	return types.COIN_ETH
}

func NewEthWallet(conf *types.Claws, ctx context.Context, conn *chain.EthConn, updateCh chan interface{}, gasPrice *chain.EthGas) *ethWallet {
	res := ethWallet{
		conf:     conf,
		ctx:      ctx,
		conn:     conn,
		updateCh: updateCh,
		gasPrice: gasPrice,
	}
	return &res
}

func (e *ethWallet) UnfoldTxs(ctx context.Context, num *big.Int) (res []types.TXN, err error) {
	b, err := e.conn.Conn.BlockByNumber(ctx, num)
	if err != nil {
		return nil, err
	}
	txs := b.Transactions()
	for _, v := range txs {
		f := strings.ToLower(v.GetFromUnsafe().String())
		txn := &ethTXN{
			from:   f,
			Hash:   v.Hash().String(),
			fee:    new(big.Int).Mul(new(big.Int).SetUint64(v.Gas()), v.GasPrice()),
			amount: v.Value(),
		}
		if v.To() != nil {
			txn.to = strings.ToLower(v.To().String())

			res = append(res, txn)
		} else {
			txn.to = ""
			res = append(res, txn)
		}
	}

	return
}

func (e *ethWallet) NotifyHead(ctx context.Context, f func(num *big.Int)) (err error) {
	ch := make(chan *ethTypes.Header)
	// cancel is an inner trigger
	ctxIn, cancel := context.WithCancel(context.TODO())
	e.once.Do(func() {
		_, err = e.conn.Conn.SubscribeNewHead(ctxIn, ch)
		if err != nil {
			return
		}
		for {
			select {
			case head := <-ch:
				f(head.Number)
			case <-e.updateCh:
				// updateCh pass a signal for reconnect
				ctxIn = context.TODO()
				ctxIn, cancel = context.WithCancel(context.TODO())
				_, err = e.conn.Conn.SubscribeNewHead(ctxIn, ch)
				if err != nil {
					return
				}
			case <-ctx.Done():
				// done all
				cancel()
			}

		}
	})
	return
}

func (e *ethWallet) Info() (info *types.Info) {
	return &types.Info{}
}

// Send send txn using bundle built for given token
func (e *ethWallet) Send(ctx context.Context, from, to types.Bundle, amount string, option *types.Option) (receipt types.Transaction, err error) {
	gasP := new(big.Int)
	if e.gasPrice != nil && e.gasPrice.GasPrice != "" {
		gasP.SetString(e.gasPrice.GasPrice, 10)
	} else {
		// 1 Gwei
		gasP.SetString("1000000000", 10)
	}

	toAddress := common.HexToAddress(to.AddressStr())

	amountBig := new(big.Int)
	amountBig.SetString(amount, 10)

	choice := &big.Int{}
	// price
	choice.SetString(option.FeeConfig, 10)
	gwei := &big.Int{}
	gwei.SetString("1000000000", 10)
	// > 1Gwei
	if choice.CmpAbs(gwei) >= 0 {
		gasP = choice
	}
	// make tx 10w for txn
	tx := ethTypes.NewTransaction(option.Nonce, toAddress, amountBig, 100000, gasP, nil)

	sb, err := hex.DecodeString(option.Secret)
	if err != nil {
		return nil, err
	}
	psk, err := crypto.ToECDSA(sb)
	if err != nil {
		return nil, err
	}
	signTx, err := ethTypes.SignTx(tx, &ethTypes.HomesteadSigner{}, psk)
	if err != nil {
		return nil, err
	}
	fmt.Println(signTx.Hash())
	if err != nil {
		fmt.Println(err)
	}
	return &ethCommit{
		signTx: signTx,
		conn:   e.conn,
	}, nil
}
