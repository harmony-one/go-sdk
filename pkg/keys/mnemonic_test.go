package keys

import (
	"testing"
)

const (
	phrase     = "crouch embrace tree can horn decrease until boil ice edit eagle chimney"
	index      = 0
	publicKey  = "0x030b64624e60c6e6758711fdf00f7e873f40a22e647a6918ba73807fab194d09ba"
	privateKey = "0xda4bc68857640103942ce7dd22a9fcdb96f3cfe0380254e81352a94ac8262ed2"
)

func TestMnemonic(t *testing.T) {
	private, public := FromMnemonicSeedAndPassphrase(phrase, index)
	sk, pkCompressed := func() (string, string) {
		dump := EncodeHex(private, public)
		return dump.PrivateKey, dump.PublicKeyCompressed
	}()
	if sk != privateKey {
		t.Errorf("Private key mismatch %s != %s", sk, privateKey)
	}
	if pkCompressed != publicKey {
		t.Errorf("Public Compressed key mismatch %s != %s", pkCompressed, publicKey)
	}
}
