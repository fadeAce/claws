package chain

import "github.com/ethereum/go-ethereum/ethclient"

type Ethereum struct {
	Url    string
	Client *EthereumClient
}

type EthereumClient struct {
	Conn *ethclient.Client
}
