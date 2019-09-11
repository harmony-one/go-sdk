package keys

import (
	"testing"
)

const (
	testMnemonic            = "quick parade stay hockey build token access sentence choice supply creek twelve"
	testPassphrase          = "edgar"
	testPrivateKey          = "0xc83827745d4e2e74e0b996b2a31ebff7bde72790bf23db839668223c07fb0299"
	testPublicCompressedKey = "0x038044dbdd8e6f901d26e3e705570b87fa6ca412f435979c60e90539a5b72b5a9f"
	testPublicKey           = "0x048044dbdd8e6f901d26e3e705570b87fa6ca412f435979c60e90539a5b72b5a9fba25f768b09c869868a5acdb9e64a2f46a3c70c04ccfb1f757f1a9ac2106418f"
)

func TestMnemonic(t *testing.T) {
	private, public := NewKeysFromMnemonicSeedAndPassphrase(testMnemonic, testPassphrase)
	sk, pkCompressed, pk := func() (string, string, string) {
		dump := EncodeHex(private, public)
		return dump.PrivateKey, dump.PublicKeyCompressed, dump.PublicKey
	}()
	if sk != testPrivateKey {
		t.Errorf("Private key mismatch %s != %s", sk, testPrivateKey)
	}
	if pkCompressed != testPublicCompressedKey {
		t.Errorf("Public Compressed key mismatch %s != %s", pkCompressed, testPublicCompressedKey)
	}
	if pk != testPublicKey {
		t.Errorf("Public Compressed key mismatch %s != %s", pk, testPublicKey)
	}
}
