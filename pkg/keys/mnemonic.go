package keys

import (
	"fmt"

	secp256k1 "github.com/btcsuite/btcd/btcec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/hd"
	"github.com/tyler-smith/go-bip39"
)

// FromMnemonicSeedAndPassphrase mimics the Harmony JS sdk in deriving the
// private, public key pair from the mnemonic, its index, and empty string password.
// Note that an index k would be the k-th key generated using the same mnemonic.
func FromMnemonicSeedAndPassphrase(mnemonic string, index int) (*secp256k1.PrivateKey, *secp256k1.PublicKey) {
	seed := bip39.NewSeed(mnemonic, "")
	master, ch := hd.ComputeMastersFromSeed(seed)
	private, _ := hd.DerivePrivateKeyForPath(
		master,
		ch,
		fmt.Sprintf("44'/1023'/0'/0/%d", index),
	)

	return secp256k1.PrivKeyFromBytes(secp256k1.S256(), private[:])
}
