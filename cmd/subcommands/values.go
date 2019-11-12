package cmd

import (
	"fmt"

	color "github.com/fatih/color"
)

const (
	hmyDocsDir      = "hmy-docs"
	defaultNodeAddr = "http://localhost:9500"
)

var (
	g           = color.New(color.FgGreen).SprintFunc()
	cookbookDoc = fmt.Sprintf(`
Cookbook of usage

note: 

1) Every subcommand recognizes a '--help' flag
2) These examples use shard 1 of testnet as argument for --node

%s
hmy --node="https://api.s1.b.hmny.io/" balance <SOME_ONE_ADDRESS>

%s
hmy --node="https://api.s1.b.hmny.io" blockchain transaction-by-hash <SOME_TX_HASH>

%s
hmy keys list

%s
hmy --node="https://api.s0.b.hmny.io/" transfer \
    --from one1yc06ghr2p8xnl2380kpfayweguuhxdtupkhqzw \
    --to one1q6gkzcap0uruuu8r6sldxuu47pd4ww9w9t7tg6 \
    --from-shard 0 --to-shard 1 --amount 200

%s
hmy --node="https://api.s0.b.hmny.io" blockchain transaction-receipt <SOME_TX_HASH>

%s
hmy keys import-ks <SOME_ABSOLUTE_PATH_TO_KEYSTORE_JSON>.key

%s
hmy keys import-private-key <secp256k1_PRIVATE_KEY>

%s
hmy keys export-private-key <ACCOUNT_ADDRESS> --passphrase harmony-one

%s
hmy keys generate-bls-key --bls-file-path /tmp/file.key

%s
hmy staking create-validator --amount 10 --validator-addr one103q7qe5t2505lypvltkqtddaef5tzfxwsse4z7 \
    --bls-pubkeys 678ec9670899bf6af85b877058bea4fc1301a5a3a376987e826e3ca150b80e3eaadffedad0fedfa111576fa76ded980c \
        6757ebfbbc53a167e4c069cdf2beecd1428316306145c8d6d97c6c6babb3ec34e7003b3db8ccfc7d79a412aec7c68c97 \
    --identity foo --details bar --name baz --max-change-rate 10 --max-rate 10 --max-total-delegation 10 \ 
    --min-self-delegation 10 --rate 10 --security-contact Leo  --website harmony.one --passphrase=''

`,
		g("1.  Check Balances"),
		g("2.  Check completed transaction"),
		g("3.  List local keys"),
		g("4.  Sending a transaction (add --wait-for-confirm=10 to wait 10 seconds for confirmation)"),
		g("5.  Check a completed transaction receipt"),
		g("6.  Import an existing keystore file"),
		g("7.  Import a keystore file using a secp256k1 private key"),
		g("8.  Export a keystore file's secp256k1 private key"),
		g("9.  Generate a BLS key then encrypt and save the private key to the specified location"),
		g("10. Create a new validator with a list of BLS keys"),
	)
)
