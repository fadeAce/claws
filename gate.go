package claws

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"

	"github.com/fadeAce/claws/line"
	"github.com/fadeAce/claws/types"
	"math/big"
	"strings"
	"sync"
	"time"
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
func SetupGate(conf *types.Claws, wildcard map[string]WalletBuilder) {
	rep := make(map[string]WalletBuilder)
	for _, v := range conf.Coins {
		var builder WalletBuilder
		//rep[v.CoinType] =
		if target, exist := wildcard[v.CoinType]; exist {
			builder = target
		} else {
			builder = setupGate(v.CoinType, conf)
		}
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
	client *line.EthConn

	url string

	gasprice *line.EthGas

	// todo: cancel of ctx or disconnection would trigger a recover of network
	ctx context.Context

	notiCh chan interface{}
}

func (e *ethBuilder) Build() Wallet {
	ctx := e.ctx
	ethWallet := line.NewEthWallet(e.config, ctx, e.client, e.notiCh, e.gasprice)
	return ethWallet
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
		ebuilder := &ethBuilder{config: conf, notiCh: make(chan interface{}, 1), gasprice: &line.EthGas{}}
		cli, err := ethclient.Dial(f(types.COIN_ETH))
		ctx := context.TODO()
		ebuilder.ctx = ctx
		if err != nil {
			fmt.Printf("create new ethereum rpc client err:%s\n", err.Error())
		}
		go func() {
			for {
				cli = ebuilder.client.Conn
				var gas *big.Int
				if cli != nil {
					gas, err = cli.SuggestGasPrice(ctx)
					if err != nil {
						fmt.Println(err.Error())
						if strings.Contains(err.Error(), "closed") ||
							strings.Contains(err.Error(), "network connection") {
							// reconnect here
							cli, err = ethclient.Dial(f(types.COIN_ETH))
							if err != nil {
								log.Error("error when reconnect source URL")
							} else {
								go func() {
									ebuilder.notiCh <- struct{}{}
								}()
								log.Info("reconnected client")
							}
							ebuilder.client.Conn = cli
						}
					}
				} else {
					cli, err = ethclient.Dial(f(types.COIN_ETH))
					if err != nil {
						log.Error("error when reconnect source URL")
					} else {
						go func() {
							ebuilder.notiCh <- struct{}{}
						}()
						log.Info("reconnected client")
					}
					ebuilder.client.Conn = cli
					time.Sleep(20 * time.Second)
					continue
				}
				time.Sleep(20 * time.Second)
				var gasStr string
				if gas != nil {
					gasStr = gas.String()
				}
				if (ebuilder.gasprice == nil || ebuilder.gasprice.GasPrice != gasStr) && gas != nil {
					ebuilder.gasprice.GasPrice = gasStr
					log.Info("updated eth gas price to ", gasStr)
				}
			}
		}()

		ebuilder.client = &line.EthConn{cli}
		return ebuilder
	}
	return nil
}
