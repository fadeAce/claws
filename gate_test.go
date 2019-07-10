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
	SetupGate(&conf, nil)
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
	SetupGate(&conf, nil)
	wallet := Builder.BuildWallet("eth")

	wallet.NotifyHead(conf.Ctx, func(num *big.Int) {
		fmt.Println(num)
	})
}

func TestSend(t *testing.T) {
	cfg, err := ioutil.ReadFile("./claws.yml")
	if cfg == nil || err != nil {
		panic("shut down with no configuration")
		return
	}
	var conf types.Claws
	err = yaml.Unmarshal(cfg, &conf)
	// first of all setup gate
	SetupGate(&conf, nil)
	wallet := Builder.BuildWallet("eth")

	_, err = wallet.Send(
		conf.Ctx,
		wallet.BuildBundle("", "", "0xf03a492fa3ce79d99b9613add1017448a83810f1"),
		wallet.BuildBundle("", "", "0x4728489Fb5c35A614c4c19450B5f964E8D794075"),
		"3000000000000000",
		&types.Option{
			Nonce:  34,
			Secret: "7b9f448ae05200d686cb982bae477e174d34c72c04d0a7464aa0d987a53d37e4",
		},
	)
	if err != nil {
		fmt.Println(err)
	}
}