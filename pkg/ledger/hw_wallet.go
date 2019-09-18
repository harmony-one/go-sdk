package ledger

import (
	"fmt"
	"sync"
	"github.com/pkg/errors"
	"golang.org/x/crypto/sha3"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/harmony-one/go-sdk/pkg/address"
	"github.com/harmony-one/harmony/core/types"
)

var (
	nanos    *NanoS        //singleton
	once     sync.Once
)

func getLedger() (*NanoS) {
	once.Do(func() {
		var err error
		nanos, err = OpenNanoS()
		if err != nil {
			log.Fatalln("Couldn't open device:", err)
			os.Exit(-1)
		}
	})

	return nanos
}

//ProcessAddressCommand list the address associated with Ledger Nano S
func GetAddress() string {
	n := getLedger()
	oneAddr, err := n.GetAddress()
	if err != nil {
		log.Fatalln("Couldn't get one address:", err)
		os.Exit(-1)
	}

	return oneAddr
}

//ProcessAddressCommand list the address associated with Ledger Nano S
func ProcessAddressCommand() {
 	n := getLedger()
	oneAddr, err := n.GetAddress()
	if err != nil {
		log.Fatalln("Couldn't get one address:", err)
		os.Exit(-1)
	}

	fmt.Printf("%-24s\t\t%23s\n", "NAME", "ADDRESS")
	fmt.Printf("%-48s\t%s\n", "Ledger Nano S",  oneAddr)
}

// SignTx signs the given transaction with the requested account.
func SignTx(tx *types.Transaction, chainID *big.Int) ([]byte, string, error) {
	var rlpEncodedTx []byte

	// Depending on the presence of the chain ID, sign with EIP155 or frontier
	if chainID != nil {
		rlpEncodedTx, _ = rlp.EncodeToBytes(
			[]interface{}{
				tx.Nonce(),
				tx.GasPrice(),
				tx.Gas(),
				tx.ShardID(),
				tx.ToShardID(),
				tx.To(),
				tx.Value(),
				tx.Data(),
				chainID, uint(0), uint(0),
			} )
	} else {
		rlpEncodedTx, _ = rlp.EncodeToBytes(
			[]interface{}{
				tx.Nonce(),
				tx.GasPrice(),
				tx.Gas(),
				tx.ShardID(),
				tx.ToShardID(),
				tx.To(),
				tx.Value(),
				tx.Data(),
			} )
	}

	n := getLedger()
	sig, err := n.SignTxn(rlpEncodedTx)
	if err != nil {
		log.Println("Couldn't sign transaction, error:", err)
		return nil, "", err
	}

	var hashBytes [32]byte
	hw := sha3.NewLegacyKeccak256()
	hw.Write(rlpEncodedTx[:])
	hw.Sum(hashBytes[:0])

	pubkey, err := crypto.Ecrecover(hashBytes[:], sig[:])
	if err != nil {
		log.Println("Ecrecover failed :", err)
		return nil, "", err
	}

	if len(pubkey) == 0 || pubkey[0] != 4 {
		log.Println("invalid public key")
		return nil, "", err
	}

	pubBytes := crypto.Keccak256(pubkey[1:65])[12:]
	signerAddr, _ := address.ConvertAndEncode("one", pubBytes)

	var r, s, v *big.Int
	if chainID != nil {
		r, s, v, err = eip155SignerSignatureValues(chainID, sig[:])
	} else {
		r, s, v, err = frontierSignatureValues(sig[:])
	}

	if err != nil {
		log.Println(err)
		return nil, "", err
	}

	// Depending on the presence of the chain ID, sign with EIP155 or frontier
	rawTx, err :=  rlp.EncodeToBytes(
		[]interface{}{
			tx.Nonce(),
			tx.GasPrice(),
			tx.Gas(),
			tx.ShardID(),
			tx.ToShardID(),
			tx.To(),
			tx.Value(),
			tx.Data(),
			v,
			r,
			s,
		} )

	return rawTx, signerAddr, err
}

func frontierSignatureValues(sig []byte) (r, s, v *big.Int, err error) {
	if len(sig) != 65 {
		return nil, nil, nil, errors.New("get signature with wrong size  from ledger nano")
	}
	r = new(big.Int).SetBytes(sig[:32])
	s = new(big.Int).SetBytes(sig[32:64])
	v = new(big.Int).SetBytes([]byte{sig[64] + 27})
	return r, s, v, nil
}

func eip155SignerSignatureValues(chainID *big.Int, sig []byte) (R, S, V *big.Int, err error) {
	R, S, V, err = frontierSignatureValues(sig)
	if err != nil {
		return nil, nil, nil, err
	}

	chainIDMul := new(big.Int).Mul(chainID, big.NewInt(2))
	if chainID.Sign() != 0 {
		V = big.NewInt(int64(sig[64] + 35))
		V.Add(V, chainIDMul)
	}
	return R, S, V, nil
}
