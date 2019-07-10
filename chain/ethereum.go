package chain

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fadeAce/claws/types"
)

type Ethereum struct {
	Conf   *types.EthConf
	Client *EthConn

	gasprice *EthGas

	// notiCh notify all updates
	notiCh chan interface{}
}


type EthConn struct {
	Conn *ethclient.Client
}

type EthGas struct {
	GasPrice string
}