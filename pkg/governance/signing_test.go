package governance

import (
	"encoding/hex"
	"os"
	"path"
	"testing"

	"github.com/btcsuite/btcd/btcec"
	"github.com/harmony-one/go-sdk/pkg/common"
	"github.com/harmony-one/harmony/accounts"
)

func TestSigning(t *testing.T) {
	// make the EIP712 data structure
	vote := Vote{
		Space:        "yam.eth",
		Proposal:     "0x21ea31e896ec5b5a49a3653e51e787ee834aaf953263144ab936ed756f36609f",
		ProposalType: "single-choice",
		Choice:       "1",
		App:          "my-app",
		Timestamp:    1660909056,
		From:         "0x9E713963a92c02317A681b9bB3065a8249DE124F",
	}
	typedData, err := vote.ToEIP712()
	if err != nil {
		t.Fatal(err)
	}

	// add a temporary key store with the below private key
	location := path.Join(os.TempDir(), "hmy-test")
	keyStore := common.KeyStoreForPath(location)
	privateKeyBytes, _ := hex.DecodeString("91c8360c4cb4b5fac45513a7213f31d4e4a7bfcb4630e9fbf074f42a203ac0b9")
	sk, _ := btcec.PrivKeyFromBytes(btcec.S256(), privateKeyBytes)
	passphrase := ""
	keyStore.ImportECDSA(sk.ToECDSA(), passphrase)
	keyStore.Unlock(accounts.Account{Address: keyStore.Accounts()[0].Address}, passphrase)
	account := accounts.Account{
		Address: keyStore.Accounts()[0].Address,
	}

	sign, err := signTypedData(keyStore, account, typedData)
	if err != nil {
		t.Fatal(err)
	}
	expectedSig := "0x6b572bacbb44efe75cad5b938a5d4fe64c3495bec28807e78989e3159e11d21d5d5568ffabcf274194830b6cb375355af423995afc7ee290e4a632b12bdbe0cc1b"
	if sign != expectedSig {
		t.Errorf("invalid sig: got %s but expected %s", sign, expectedSig)
	}

	os.RemoveAll(location)
}

// The below NodeJS code was used to generate the above signature
// import snapshot from '@snapshot-labs/snapshot.js';
// import { Wallet } from "ethers";
// const hub = 'https://hub.snapshot.org';
// const client = new snapshot.Client712(hub);
// const wallet = new Wallet("91c8360c4cb4b5fac45513a7213f31d4e4a7bfcb4630e9fbf074f42a203ac0b9");
// const receipt = await client.vote(wallet, await wallet.getAddress(), {
// space: 'yam.eth',
// proposal: '0x21ea31e896ec5b5a49a3653e51e787ee834aaf953263144ab936ed756f36609f',
// type: 'single-choice',
// choice: 1,
// app: 'my-app',
// timestamp: 1660909056,
// });

// package.json
// "@snapshot-labs/snapshot.js": "^0.4.18"
// "ethers": "^5.6.9"
