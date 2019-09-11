package keys

import (
	"fmt"

	secp256k1 "github.com/btcsuite/btcd/btcec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/hd"
	"github.com/tyler-smith/go-bip39"
)

// const (
// 	TestMnemonic   = "quick parade stay hockey build token access sentence choice supply creek twelve"
// 	TestPassphrase = "edgar"
// 	TestPrivateKey = "0xc83827745d4e2e74e0b996b2a31ebff7bde72790bf23db839668223c07fb0299"
// 	TestPublicKey  = "0x038044dbdd8e6f901d26e3e705570b87fa6ca412f435979c60e90539a5b72b5a9f"
// )

func FromMnemonicSeedAndPassphrase(mnemonic, passphrase string) (*secp256k1.PrivateKey, *secp256k1.PublicKey) {
	seed := bip39.NewSeed(mnemonic, passphrase)
	master, ch := hd.ComputeMastersFromSeed(seed)
	// TODO Come back to idea of index/issue

	index := 0
	private, _ := hd.DerivePrivateKeyForPath(
		master,
		ch,
		fmt.Sprintf("44'/1023'/0'/0/%d", index),
	)

	return secp256k1.PrivKeyFromBytes(secp256k1.S256(), private[:])

	// p1 := publicK.SerializeCompressed()
	// p2 := privateK.Serialize()

	// ePublic := hexutil.Encode(p1)
	// ePrivate := hexutil.Encode(p2)

	// fmt.Printf("Public: %s \nPrivate: %s\nMnemonic: %s\n",
	// 	ePublic,
	// 	ePrivate,
	// 	TestMnemonic,
	// )

	// if ePublic != TestPublicKey {
	// 	fmt.Printf("Error, not public matching %s\n", TestPublicKey)
	// }
	// if ePrivate != TestPrivateKey {
	// 	fmt.Printf("Error, not private matching %s\n", TestPrivateKey)
	// }
	// return p1, p2, nil
	// account, error :=
	// 	DefaultKS.ImportECDSA(privateK.ToECDSA(), passphrase)

	// fmt.Println(account, error)
}

// func T() {
// 	NewAccountFromMnemonicSeedAndPassphrase(TestMnemonic, TestPassphrase)
// }
