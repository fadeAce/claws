package claws

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fadeAce/claws/line"
	"github.com/fadeAce/claws/types"

	"sync"
)

// this gate is for matching cgate at claws
// it's eternal and not to be modified
var Builder = func() *builder {
	return &builder{
		builder: make(map[string]WalletBuilder),
		mu:      sync.RWMutex{},
	}
}()

type builder struct {
	builder map[string]WalletBuilder
	mu      sync.RWMutex
}

type WalletBuilder interface {
	Build() Wallet
}

// builder build wallet in different type given by params
func (b *builder) BuildWallet(typ string) Wallet {
	b.mu.Lock()
	defer b.mu.Unlock()
	builder := b.builder[typ]
	if builder == nil {
		return nil
	}
	return builder.Build()
}

// this is based on cfg itself, it would load all config to make intersection with code pre-define
func SetupGate(conf *types.Claws) {
	rep := make(map[string]WalletBuilder)
	for _, v := range conf.Coins {
		//rep[v.CoinType] =
		builder := setupGate(v.CoinType, conf)
		if builder != nil {
			rep[v.CoinType] = builder
		}
	}
	// make up builder map with a brand new copy
	Builder.builder = rep
	conf.Ctx = context.TODO()
}

// ethBuilder holds connection to eth nodes, if it's disconnected a recover would go automatically
type ethBuilder struct {
	config *types.Claws

	// connection to nodes
	client *ethclient.Client
	// todo: cancel of ctx or disconnection would trigger a recover of network
	ctx context.Context
}

func (e *ethBuilder) Build() Wallet {
	ctx := e.ctx
	ethWallet := line.NewEthWallet(e.config, ctx, e.client)
	return &ethWallet
}

// there is a list of coins to be initiated which make up the settings
func setupGate(typ string, conf *types.Claws) WalletBuilder {
	// f returns the URL represented RPC node
	f := func(typ string) string {
		for _, v := range conf.Coins {
			if v.CoinType == typ {
				return v.Url
			}
		}
		return ""
	}
	switch typ {
	case types.COIN_BTC:
		return nil
	case types.COIN_ETH:
		ebuilder := &ethBuilder{config: conf}
		cli, err := ethclient.Dial(f(types.COIN_ETH))
		if err != nil {
			fmt.Printf("create new ethereum rpc client err:%s\n", err.Error())
		} else {
			fmt.Println("create new ethereum rpc client success")
		}
		etx := context.TODO()
		ebuilder.ctx = etx
		ebuilder.client = cli
		return ebuilder
	}
	return nil
}
