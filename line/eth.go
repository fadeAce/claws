package line

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	types2 "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/marblebank/claws"
	"github.com/marblebank/claws/addr"
	"github.com/marblebank/config"
	"github.com/marblebank/types"
	"github.com/opentracing/opentracing-go/log"
)

type ethWallet struct {
	conf *config.Marble
	ctx  context.Context
	conn *ethclient.Client
}

type ethBundle struct {
	pub, prv, add string
}

type ethTXN struct {
	Hash string
}

func (etxn *ethTXN) HexStr() string {
	return ""
}

func (etxn *ethTXN) SetStr(res string) {
	etxn.Hash = res
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
func (e *ethWallet) NewAddr() claws.Bundle {
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
) claws.Bundle {
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
) claws.TXN {
	return &ethTXN{
		Hash: tx,
	}
}

func (e *ethWallet) Info() *types.Info {
	// filter config for wallet
	for _, v := range e.conf.Coins {
		info := &types.Info{}
		if e.Type() == v.CoinType {
			info.Fee = v.Fee
			info.FeeCoin = v.FeeCoin
			info.Name = v.CoinType
			info.Decimal = v.Decimal
			info.DisplayShort = v.ShortName
			info.Display = v.FullName
		}
		return info
	}
	return nil
}

// withdraw to certain place
func (e *ethWallet) Withdraw(addr claws.Bundle) *types.TxnInfo {
	return nil
}

// seek for tx , keep track it
func (e *ethWallet) Seek(txn claws.TXN) bool {
	// seek the txn hash of it
	hash := txn.HexStr()
	hs := common.HexToHash(hash)
	reciept, err := e.conn.TransactionReceipt(e.ctx, hs)
	if err != nil {
		log.Error(err)
		return false
	}
	if reciept.Status == types2.ReceiptStatusSuccessful {
		return true
	}
	return false
}

// seek for tx , keep track it
func (e *ethWallet) Balance(bundle claws.Bundle) (string, error) {
	add := bundle.AddressStr()
	address := common.HexToAddress(add)
	balance, err := e.conn.BalanceAt(e.ctx, address, nil)
	return balance.String(), err
}

// seek for tx , keep track it
func (e *ethWallet) Type() string {
	return types.COIN_ETH
}

func NewEthWallet(conf *config.Marble, ctx context.Context, conn *ethclient.Client) ethWallet {
	res := ethWallet{
		conf: conf,
		ctx:  ctx,
		conn: conn,
	}
	return res
}
