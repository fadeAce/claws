package addr

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/marblebank/types"
	"github.com/pborman/uuid"

	crand "crypto/rand"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/crypto"
)

// automatic decode from type if failed it would trigger panic
type Addr interface {
	GetType() string
	HexPubStr() string
	HexPrvStr() string
	HexAddrStr() string
}

type EthAddr struct {
	typ string
	key *keystore.Key
}

func (e *EthAddr) HexPubStr() string {
	return ""
}
func (e *EthAddr) HexPrvStr() string {
	return hex.EncodeToString(e.key.PrivateKey.D.Bytes())
}
func (e *EthAddr) HexAddrStr() string {
	return e.key.Address.Hex()
}

func (e *EthAddr) GetType() string {
	return types.COIN_ETH
}

func GenerateAddr(typ string) (Addr, error) {
	switch typ {
	case "btc":
		return btcAddr()
	case "eth":
		return ethAddr()
	}
	return nil, errors.New("can't generate " + typ + " addr")
}

func btcAddr() (Addr, error) {
	prv, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {

	}
	pub := prv.PubKey()
	pbk := pub.SerializeCompressed()
	pkHash := btcutil.Hash160(pbk)
	ad, err := btcutil.NewAddressPubKeyHash(pkHash, &chaincfg.Params{
		// this stands for test net
		PubKeyHashAddrID: 0x6f,
		// this stands for main net
		//PubKeyHashAddrID: 0x00,
	})
	fmt.Println(ad.String())
	fmt.Println(hex.EncodeToString(prv.Serialize()))
	return nil, nil
}

func ethAddr() (Addr, error) {
	privateKeyECDSA, err := ecdsa.GenerateKey(crypto.S256(), crand.Reader)
	if err != nil {
		return nil, err
	}
	id := uuid.NewRandom()
	key := &keystore.Key{
		Id:         id,
		Address:    crypto.PubkeyToAddress(privateKeyECDSA.PublicKey),
		PrivateKey: privateKeyECDSA,
	}
	ad := &EthAddr{
		typ: "eth",
		key: key,
	}
	return ad, nil
}
