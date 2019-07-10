package chain

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fadeAce/claws/types"
)

type Ethereum struct {
	Conf   *types.EthConf
	Client *EthereumClient
}

type EthereumClient struct {
	Conn *ethclient.Client
}
