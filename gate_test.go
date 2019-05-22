package claws

import (
	"fmt"
	"github.com/fadeAce/claws/types"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"math/big"
	"testing"
)

func TestWalletBase(t *testing.T) {
	cfg, err := ioutil.ReadFile("./claws.yml")
	if cfg == nil || err != nil {
		panic("shut down with no configuration")
		return
	}
	var conf types.Claws
	err = yaml.Unmarshal(cfg, &conf)
	// first of all setup gate
	SetupGate(&conf)
	wallet := Builder.BuildWallet("eth")
	b := wallet.NewAddr()
	fmt.Println(b)

	//num := big.NewInt(4419795)
	txns, err := wallet.UnfoldTxs(conf.Ctx, big.NewInt(4356126))
	for _, v := range txns {
		fmt.Println("from ", v.FromStr(), " to ", v.ToStr(), " hash ", v.HexStr())
		fmt.Println(" fee ", v.FeeStr(), " amount ", v.AmountStr())
	}
}

func TestEventBase(t *testing.T) {
	cfg, err := ioutil.ReadFile("./claws.yml")
	if cfg == nil || err != nil {
		panic("shut down with no configuration")
		return
	}
	var conf types.Claws
	err = yaml.Unmarshal(cfg, &conf)
	// first of all setup gate
	SetupGate(&conf)
	wallet := Builder.BuildWallet("eth")

	wallet.NotifyHead(conf.Ctx, func(num *big.Int) {
		fmt.Println(num)
	})
}
