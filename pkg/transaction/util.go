package transaction

import (
	"encoding/base64"
	"encoding/hex"
	"math/big"

	"github.com/btcsuite/btcd/btcec"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/harmony-one/go-sdk/pkg/address"
	"github.com/harmony-one/go-sdk/pkg/common"
	"github.com/harmony-one/go-sdk/pkg/sharding"
	"github.com/harmony-one/go-sdk/pkg/store"
	"github.com/harmony-one/go-sdk/pkg/validation"
	"github.com/harmony-one/harmony/common/denominations"
	"github.com/harmony-one/harmony/core"
)

const (
	node = "http://localhost:9500"
)

func SignTx(tx *Transaction, privateKey, passphrase, accountName string, chainID *big.Int) ([]byte, error) {
	// create local account using the imported privat key and passphrase
	if privateKey[:2] == "0x" {
		privateKey = privateKey[2:]
	}
	privateKeyBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		return nil, err
	}
	if len(privateKeyBytes) != common.Secp256k1PrivateKeyBytesLength {
		return nil, common.ErrBadKeyLength
	}
	// btcec.PrivKeyFromBytes only returns a secret key and public key
	sk, _ := btcec.PrivKeyFromBytes(btcec.S256(), privateKeyBytes)
	ks := store.FromAccountName(accountName)
	account, err := ks.ImportECDSA(sk.ToECDSA(), passphrase)
	if err != nil {
		return nil, err
	}

	from := address.ToBech32(account.Address)
	ks, acc, err := store.UnlockedKeystore(from, passphrase)
	if err != nil {
		return nil, err
	}

	signedTx, err := ks.SignTx(*acc, tx, chainID)
	if err != nil {
		return nil, err
	}

	enc, _ := rlp.EncodeToBytes(signedTx)
	hexSignature := hexutil.Encode(enc)

	return []byte(hexSignature), nil
}

func NewTx(nonce, gasLimit uint64,
	to string,
	shardID, toShardID uint32,
	amount, gasPrice *big.Int,
	data []byte) (*Transaction, error) {

	// verify the shard ids
	s, err := sharding.Structure(node)
	if err != nil {
		return nil, err
	}
	err = validation.ValidShardIDs(shardID, toShardID, uint32(len(s)))
	if err != nil {
		return nil, err
	}

	// set amount
	nanoVal := amount.Mul(amount, big.NewInt(denominations.Nano))
	amountBigInt := big.NewInt(nanoVal.Int64())
	amt := amountBigInt.Mul(amountBigInt, big.NewInt(denominations.Nano))

	// set gas limit
	inputData, _ := base64.StdEncoding.DecodeString(string(data))
	gas, _ := core.IntrinsicGas(inputData, false, true)

	return NewTransaction(nonce, gas, address.Parse(to), shardID, toShardID, amt, nil, data), nil
}
