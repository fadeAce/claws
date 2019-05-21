package addr

import (
	"context"
	types2 "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"os"

	"bufio"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/mobile"
	"github.com/fadeAce/claws/types"
	"log"
	"strings"
	"testing"
)

func TestGenerateBtcAddr(t *testing.T) {
	add, err := GenerateAddr("btc")
	if err != nil {
		t.Error(err)
	}
	fmt.Println("add is ", add)

}

func TestGenerateEthAddr(t *testing.T) {
	add, err := GenerateAddr("eth")
	if err != nil {
		t.Error(err)
	}
	key := add.(*EthAddr)
	fmt.Println("add is ", key.key.Address.String(), " ", key.key.PrivateKey.D)
}

func TestEthKeyFile(t *testing.T) {
	add, err := GenerateAddr("eth")
	if err != nil {
		t.Error(err)
	}
	key := add.(*EthAddr)
	ks := keystore.NewKeyStorePlain("/Users/terrill/main/GoProject/marblebank/key.dat")
	err = ks.StoreKey("/Users/terrill/main/GoProject/marblebank/key.dat", key.key, "1234")
	if err != nil {
		t.Error(err)
	}
}

func TestEthKeyFileRestore(t *testing.T) {
	ks := keystore.NewKeyStorePlain("/Users/terrill/main/GoProject/marblebank/key.dat")
	address := common.HexToAddress("0xf03a492fa3ce79d99b9613add1017448a83810f1")
	key, err := ks.GetKey(address, "/Users/terrill/main/GoProject/marblebank/key.dat", "")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(key.PrivateKey.D)
	fmt.Println(hex.EncodeToString(key.PrivateKey.D.Bytes()))
	fmt.Println(key.Address.String())
}

func TestBtcRPC(t *testing.T) {
	// Connect to local bitcoin core RPC server using HTTP POST mode.
	connCfg := &rpcclient.ConnConfig{
		Host:         "127.0.0.1:18332",
		User:         "alice",
		Pass:         "alice",
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default
	}
	// Notice the notification parameter is nil since notifications are
	// not supported in HTTP POST mode.
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Shutdown()

	// Get the current block count.
	blockCount, err := client.GetBlockCount()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Block count: %d", blockCount)
}

func TestEthRPC(t *testing.T) {
	// NewEthereumClient函数只是创建一个EthereumClient结构，并设置了HTTP连接的一些参数如的head的一些属性，并没有节点建立连接
	cli, err := geth.NewEthereumClient("http://127.0.0.1:8545")
	if err != nil {
		fmt.Printf("create new ethereum rpc client err:%s\n", err.Error())
	} else {
		fmt.Println("create new ethereum rpc client success")
	}
	eth_ctx := geth.NewContext()
	block, err2 := cli.GetBlockByNumber(eth_ctx, -1)
	fmt.Printf("ethereum mobile Context:%+v\n", eth_ctx)
	if err2 != nil {
		fmt.Printf("get block err:%s\n", err2.Error())
	} else {
		fmt.Printf("block num:%+v\n", block.GetNumber())
	}
}

func TestEthBalance(t *testing.T) {
	cli, err := geth.NewEthereumClient("http://127.0.0.1:8545")
	if err != nil {
		fmt.Printf("create new ethereum rpc client err:%s\n", err.Error())
	} else {
		fmt.Println("create new ethereum rpc client success")
	}
	eth_ctx := geth.NewContext()
	adb, _ := hex.DecodeString("f03a492fa3ce79d99b9613add1017448a83810f1")
	a, _ := geth.NewAddressFromBytes(adb)
	block, err2 := cli.GetBalanceAt(eth_ctx, a, -1)
	fmt.Printf("ethereum mobile Context:%+v\n", eth_ctx)
	if err2 != nil {
		fmt.Printf("get block err:%s\n", err2.Error())
	} else {
		fmt.Printf("block:%+v\n", block)
	}
}

func TestEthPendingTransaction(t *testing.T) {
	client, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		fmt.Println(err)
	}
	a, err := client.PendingTransactionCount(context.TODO())
	fmt.Println(a)
}

//func TestSendEth(t *testing.T) {
//	cli, err := geth.NewEthereumClient("http://127.0.0.1:8545")
//	if err != nil {
//		fmt.Printf("create new ethereum rpc client err:%s\n", err.Error())
//	} else {
//		fmt.Println("create new ethereum rpc client success")
//	}
//	eth_ctx := geth.NewContext()
//
//	ad, err := geth.NewAddressFromHex("0xFC8eFa516089C2c1Fd606914c94137a98653b676")
//	if err != nil {
//		fmt.Println(err)
//	}
//	tx := geth.NewTransaction(
//		0, ad, geth.NewBigInt(21000000), 200000, geth.NewBigInt(21000), nil)
//
//	signer := &etype.HomesteadSigner{}
//
//	a := new(etype.Address)
//	if err := a.SetHex(hex); err != nil {
//		return nil, err
//	}
//
//	txOrigin := etype.NewTransaction(uint64(nonce), to.address, amount.bigint, uint64(gasLimit), gasPrice.bigint, nil)
//
//	signature, err := crypto.Sign(signer.Hash(tx).Bytes(), key)
//	stx, err := tx.WithSignature()
//	if err != nil {
//		fmt.Println(err)
//	}
//	err = cli.SendTransaction(eth_ctx, stx)
//	if err != nil {
//		fmt.Println(err)
//	}
//}

func TestSendEth(t *testing.T) {
	client, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		fmt.Println(err)
	}

	toAddress := common.HexToAddress("FC8eFa516089C2c1Fd606914c94137a98653b676")

	// make tx
	tx := types2.NewTransaction(0, toAddress, big.NewInt(30000000000000000), 40000, big.NewInt(1000), nil)

	key := new(keystore.Key)

	// file testing

	fd, err := os.Open("key.dat")
	if err != nil {
		fmt.Println(err)
	}
	defer fd.Close()
	if err := json.NewDecoder(fd).Decode(key); err != nil {
		fmt.Println(err)
	}
	signTx, err := types2.SignTx(tx, &types2.HomesteadSigner{}, key.PrivateKey)
	fmt.Println(signTx.Hash())
	if err != nil {
		fmt.Println(err)
	}
	err = client.SendTransaction(context.TODO(), signTx)
	if err != nil {
		fmt.Println("error :", err)
	}
	fmt.Println(signTx.Hash().Hex())
}

func TestEstimate(t *testing.T) {
	cli, err := geth.NewEthereumClient("http://127.0.0.1:8545")
	if err != nil {
		fmt.Printf("create new ethereum rpc client err:%s\n", err.Error())
	} else {
		fmt.Println("create new ethereum rpc client success")
	}
	eth_ctx := geth.NewContext()

	toad, err := geth.NewAddressFromHex("FC8eFa516089C2c1Fd606914c94137a98653b676")
	frad, err := geth.NewAddressFromHex("f03a492fa3ce79d99b9613add1017448a83810f1")
	if err != nil {
		fmt.Println(err)
	}

	callmsg := geth.NewCallMsg()

	callmsg.SetData(nil)
	callmsg.SetFrom(frad)
	callmsg.SetGas(21000)
	//callmsg.SetGasPrice(geth.NewBigInt(18000000000))
	callmsg.SetGasPrice(geth.NewBigInt(200000))
	callmsg.SetTo(toad)
	callmsg.SetValue(geth.NewBigInt(21000000))

	res, err := cli.EstimateGas(eth_ctx, callmsg)
	fmt.Println("estimates gas is :", res)
	if err != nil {
		fmt.Println(err)
	}
}

func TestEthAddrInterface(t *testing.T) {
	add, err := GenerateAddr(types.COIN_ETH)
	if err != nil {
		t.Error(err)
	}
	key := add.(*EthAddr)
	fmt.Println(key.key.Id.String())
	fmt.Println(key.key.Address.String())
	fmt.Println(hex.EncodeToString(key.key.PrivateKey.D.Bytes()))

	// recover directly from here not a located based dat file
	fmt.Println("--- line")
	addressStr := strings.Replace(key.key.Address.String(), "0x", "", -1)

	k := GetKey(addressStr, key.key.Id.String(), hex.EncodeToString(key.key.PrivateKey.D.Bytes()))
	fmt.Println(k.Id.String())
	fmt.Println(k.Address.String())
	fmt.Println(hex.EncodeToString(k.PrivateKey.D.Bytes()))

}

//func TestBtcBalance(t *testing.T) {
//	// Connect to local bitcoin core RPC server using HTTP POST mode.
//	connCfg := &rpcclient.ConnConfig{
//		Host:         "127.0.0.1:18332",
//		User:         "alice",
//		Pass:         "alice",
//		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
//		DisableTLS:   true, // Bitcoin core does not provide TLS by default
//	}
//	// Notice the notification parameter is nil since notifications are
//	// not supported in HTTP POST mode.
//	client, err := rpcclient.New(connCfg, nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer client.Shutdown()
//	adStr := "mzUidJJwrTMoczCapRtzi1JDfxqjvHAt4A"
//	pubh, v, err := base58.CheckDecode(adStr)
//	ad, err := btcutil.NewAddressPubKeyHash(pubh, &chaincfg.Params{
//		PubKeyHashAddrID: v,
//	})
//	// Get the current block count.
//	fmt.Println(v)
//	balance, err := client.GetRawTransactionVerbose()
//	if err != nil {
//		log.Fatal(err)
//	}
//	log.Printf("Block count: %d", balance)
//}

type KeyJ struct {
	//Id         string `json:"id"`
	Address    string `json:"address"`
	Privatekey string `json:"privatekey"`
	V          int    `json:"version"`
}

func GetKey(a, i, p string) *keystore.Key {
	key := new(keystore.Key)
	k := KeyJ{
		//Id:         i,
		Address:    a,
		Privatekey: p,
		V:          3,
	}
	keyJ, _ := json.Marshal(k)

	//// file testing
	//ioutil.WriteFile("tmp.dat", keyJ, 0777)
	//
	//fd, err := os.Open("tmp.dat")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//defer fd.Close()
	//if err := json.NewDecoder(fd).Decode(key); err != nil {
	//	fmt.Println(err)
	//}
	rd := bytes.NewReader(keyJ)
	buf := bufio.NewReader(rd)
	_ = json.NewDecoder(buf).Decode(key)
	fmt.Println(string(keyJ))
	return key
}

func TestEthReceipt(t *testing.T) {
	client, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		fmt.Println(err)
	}
	if err != nil {
		fmt.Printf("create new ethereum rpc client err:%s\n", err.Error())
	} else {
		fmt.Println("create new ethereum rpc client success")
	}

	txHs := "0xf94f517b0ecff3d73a2f68e9020150bc98288afc8369af2e36292d8e2c4bd303"
	hs := common.HexToHash(txHs)
	r, err := client.TransactionReceipt(context.TODO(), hs)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(r)
	if r.Status == types2.ReceiptStatusSuccessful {
		fmt.Println("success !")
	}
}

func TestEthWallet(t *testing.T) {
	//cli := claws.Wallet()
}


