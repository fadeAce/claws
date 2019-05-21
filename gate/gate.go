package gate
//
//import (
//	"context"
//	"fmt"
//	"github.com/ethereum/go-ethereum/ethclient"
//	"github.com/fadeAce/claws"
//	"github.com/fadeAce/claws/line"
//	"github.com/fadeAce/claws/types"
//
//	"sync"
//)
//
//// this gate is for matching cgate at claws
//// it's eternal and not to be modified
//var Builder = func() *builder {
//	return &builder{
//		builder: make(map[string]WalletBuilder),
//		mu:      sync.RWMutex{},
//	}
//}()
//
//type builder struct {
//	builder map[string]WalletBuilder
//	mu      sync.RWMutex
//}
//
//type WalletBuilder interface {
//	Build() claws.Wallet
//}
//
//// builder build wallet in different type given by params
//func (b *builder) BuildWallet(typ string) claws.Wallet {
//	b.mu.Lock()
//	defer b.mu.Unlock()
//	builder := b.builder[typ]
//	if builder == nil {
//		return nil
//	}
//	return builder.Build()
//}
//
//// this is based on cfg itself, it would load all config to make intersection with code pre-define
//func SetupGate(conf *types.Marble) {
//	rep := make(map[string]WalletBuilder)
//	for _, v := range conf.Coins {
//		//rep[v.CoinType] =
//		builder := setupGate(v.CoinType, conf)
//		if builder != nil {
//			rep[v.CoinType] = builder
//		}
//	}
//	// make up builder map with a brand new copy
//	Builder.builder = rep
//
//}
//
//// ethBuilder holds connection to eth nodes, if it's disconnected a recover would go automatically
//type ethBuilder struct {
//	config *types.Marble
//
//	// connection to nodes
//	client *ethclient.Client
//	// todo: cancel of ctx or disconnection would trigger a recover of network
//	ctx context.Context
//}
//
//func (e *ethBuilder) Build() claws.Wallet {
//	ctx := e.ctx
//	ethWallet := line.NewEthWallet(e.config, ctx, e.client)
//	return &ethWallet
//}
//
//// there is a list of coins to be initiated which make up the settings
//func setupGate(typ string, conf *types.Marble) WalletBuilder {
//	switch typ {
//	case types.COIN_BTC:
//		return nil
//	case types.COIN_ETH:
//		ebuilder := &ethBuilder{config: conf}
//		// gate do connect to RPC nodes and share it's connection
//		//cli, err := geth.NewEthereumClient("http://127.0.0.1:8545")
//		cli, err := ethclient.Dial("http://127.0.0.1:8545")
//		if err != nil {
//			fmt.Printf("create new ethereum rpc client err:%s\n", err.Error())
//		} else {
//			fmt.Println("create new ethereum rpc client success")
//		}
//		etx := context.TODO()
//		ebuilder.ctx = etx
//		ebuilder.client = cli
//		return ebuilder
//	}
//	return nil
//}
