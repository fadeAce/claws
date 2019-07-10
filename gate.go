package claws

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fadeAce/claws/chain"
	"github.com/fadeAce/claws/line"
	"github.com/fadeAce/claws/types"
	"github.com/rs/zerolog/log"
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
	// chain configurations here
	EthChain *chain.Ethereum
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

	// setup ETH chain
	setupETH(conf.Eth)

	//setup btc .. eos .. etc

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

	if conf.Eth != nil {
		Builder.EthChain = &chain.Ethereum{
			Conf:   conf.Eth,
			Client: nil,
		}
	}

	conf.Ctx = context.TODO()
}

// ethBuilder holds connection to eth nodes, if it's disconnected a recover would go automatically
type ethBuilder struct {
	config *types.Claws

	// connection to nodes
	client *chain.EthConn

	url string

	gasprice *chain.EthGas

	// todo: cancel of ctx or disconnection would trigger a recover of network
	ctx context.Context

	notiCh chan interface{}
}

func (e *ethBuilder) Build() Wallet {
	ctx := e.ctx
	ethWallet := line.NewEthWallet(e.config, ctx, e.client, e.notiCh, e.gasprice)
	return ethWallet
}

type erc20Builder struct {
	*sync.RWMutex
	config   *types.Claws
	client   *line.Erc20Client
	url      string
	gasprice *chain.EthGas
	ctx      context.Context
	notiCh   chan interface{}
}

func (e *erc20Builder) Build() Wallet {
	ctx := e.ctx
	coinCfg := getChainConf("erc20", e.config)
	erc20Wallet := line.NewERC20Wallet(e.config, coinCfg, ctx, e.client, e.RWMutex)
	return erc20Wallet
}

func getChainConf(typ string, conf *types.Claws) *types.Coins {
	for _, v := range conf.Coins {
		if v.CoinType == typ {
			return &v
		}
	}
	panic("configure not find")
	return nil
}

func setupETH(conf *types.EthConf) {
	////ebuilder := &ethBuilder{config: conf, notiCh: make(chan interface{}, 1), gasprice: &chain.EthGas{}}
	//cli, err := ethclient.Dial(conf.Url)
	//ctx := context.TODO()
	//if err != nil {
	//	fmt.Printf("create new ethereum rpc client err:%s\n", err.Error())
	//}
	//go func() {
	//	for {
	//		cli = Builder.EthChain.Client.Conn
	//		var gas *big.Int
	//		if cli != nil {
	//			gas, err = cli.SuggestGasPrice(ctx)
	//			if err != nil {
	//				fmt.Println(err.Error())
	//				if strings.Contains(err.Error(), "closed") ||
	//					strings.Contains(err.Error(), "network connection") {
	//					// reconnect here
	//					cli, err = ethclient.Dial(conf.Url)
	//					if err != nil {
	//						log.Error().Msg("error when reconnect source URL")
	//					} else {
	//						go func() {
	//							Builder.EthChain.notiCh <- struct{}{}
	//						}()
	//						log.Info().Msg("reconnected client")
	//					}
	//					ebuilder.client.Conn = cli
	//				}
	//			}
	//		} else {
	//			cli, err = ethclient.Dial(conf.Url)
	//			if err != nil {
	//				log.Error().Msg("error when reconnect source URL")
	//			} else {
	//				go func() {
	//					ebuilder.notiCh <- struct{}{}
	//				}()
	//				log.Info().Msg("reconnected client")
	//			}
	//			ebuilder.client.Conn = cli
	//			time.Sleep(20 * time.Second)
	//			continue
	//		}
	//		time.Sleep(20 * time.Second)
	//		var gasStr string
	//		if gas != nil {
	//			gasStr = gas.String()
	//		}
	//		if (ebuilder.gasprice == nil || ebuilder.gasprice.GasPrice != gasStr) && gas != nil {
	//			ebuilder.gasprice.GasPrice = gasStr
	//			log.Info().Msgf("updated eth gas price to %s", gasStr)
	//		}
	//	}
	//}()
	//
	//ebuilder.client = &chain.EthConn{cli}
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
		ebuilder := &ethBuilder{config: conf, notiCh: make(chan interface{}, 1), gasprice: &chain.EthGas{}}
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
								log.Error().Msg("error when reconnect source URL")
							} else {
								go func() {
									ebuilder.notiCh <- struct{}{}
								}()
								log.Info().Msg("reconnected client")
							}
							ebuilder.client.Conn = cli
						}
					}
				} else {
					cli, err = ethclient.Dial(f(types.COIN_ETH))
					if err != nil {
						log.Error().Msg("error when reconnect source URL")
					} else {
						go func() {
							ebuilder.notiCh <- struct{}{}
						}()
						log.Info().Msg("reconnected client")
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
					log.Info().Msgf("updated eth gas price to %s", gasStr)
				}
			}
		}()

		ebuilder.client = &chain.EthConn{cli}
		return ebuilder
	case types.COIN_ERC20:
		return newErc20Builder(context.Background(), conf)
	}
	return nil
}

func newErc20Builder(ctx context.Context, conf *types.Claws) WalletBuilder {
	cfg := getChainConf("erc20", conf)
	erc20_builder := &erc20Builder{
		config:   conf,
		notiCh:   make(chan interface{}, 1),
		gasprice: &chain.EthGas{},
		RWMutex:  &sync.RWMutex{},
	}

	if c, err := line.NewERC20Client(ctx, cfg.Url); err == nil {
		erc20_builder.Lock()
		erc20_builder.client = c
		erc20_builder.Unlock()
	} else {
		panic(err)
	}

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			<-ticker.C
			if client, err := ethclient.Dial(cfg.Url); err != nil || erc20_builder.client.Closed {
				//erc20_builder.Lock()
				erc20_builder.client.Closed = true
				if err := erc20_builder.client.Reconnect(); err == nil {
					erc20_builder.client.Closed = false
					log.Info().Msg("ERC20 network reconnect succeed.")
				} else {
					log.Info().Msg("ERC20 network reconnect failed.")
				}
				//erc20_builder.Unlock()
			} else {
				client.Close()
			}
		}
	}()
	return erc20_builder
}
