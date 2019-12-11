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
Cookbook of Usage

Note:

1) Every subcommand recognizes a '--help' flag
2) These examples use shard 1 of testnet as argument for --node

Examples:

%s
hmy --node="https://api.s1.p.hmny.io/" balances <SOME_ONE_ADDRESS>

%s
hmy --node="https://api.s1.p.hmny.io" blockchain transaction-by-hash <SOME_TX_HASH>

%s
hmy keys list

%s
hmy --node="https://api.s1.p.hmny.io/" transfer \
    --from one1yc06ghr2p8xnl2380kpfayweguuhxdtupkhqzw \
    --to one1q6gkzcap0uruuu8r6sldxuu47pd4ww9w9t7tg6 \
    --from-shard 0 --to-shard 1 --amount 200

%s
hmy --node="https://api.s1.p.hmny.io" blockchain transaction-receipt <SOME_TX_HASH>

%s
hmy keys add --recover

%s
hmy keys import-ks <SOME_ABSOLUTE_PATH_TO_KEYSTORE_JSON>.key

%s
hmy keys import-private-key <secp256k1_PRIVATE_KEY>

%s
hmy keys export-private-key <ACCOUNT_ADDRESS> --passphrase <YOUR_PASSWORD>

%s
hmy keys generate-bls-key --bls-file-path /tmp/file.key

%s
hmy --node="https://api.s1.p.hmny.io" staking create-validator --amount 10 --validator-addr <SOME_ONE_ADDRESS> \
    --bls-pubkeys <BLS_KEY_1>,<BLS_KEY_2>,<BLS_KEY_3> \
    --identity foo --details bar --name baz --max-change-rate 0.1 --max-rate 0.1 --max-total-delegation 10 \
    --min-self-delegation 10 --rate 0.1 --security-contact Leo  --website harmony.one --passphrase <YOUR_PASSWORD>

%s
hmy --node="https://api.s1.p.hmny.io" staking edit-validator \
    --validator-addr <SOME_ONE_ADDRESS> --identity foo --details bar \
    --name baz --security-contact EK --website harmony.one \
    --min-self-delegation 0 --max-total-delegation 10 --rate 0.1\
    --add-bls-key <SOME_BLS_KEY> --remove-bls-key <OTHER_BLS_KEY> --passphrase <YOUR_PASSWORD>

%s
hmy --node="https://api.s1.p.hmny.io" staking delegate \
    --delegator-addr <SOME_ONE_ADDRESS> --validator-addr <VALIDATOR_ONE_ADDRESS> \
    --amount 10 --passphrase <YOUR_PASSWORD>

%s
hmy --node="https://api.s1.p.hmny.io" staking undelegate \
    --delegator-addr <SOME_ONE_ADDRESS> --validator-addr <VALIDATOR_ONE_ADDRESS> \
    --amount 10 --passphrase <YOUR_PASSWORD>

%s
hmy --node="https://api.s1.p.hmny.io" staking collect-rewards \
    --delegator-addr <SOME_ONE_ADDRESS> --passphrase <YOUR_PASSWORD>

%s
hmy --node="https://api.s1.p.hmny.io" blockchain validator all-active

%s
hmy failures staking

`,
		g("1.  Check account balance on given chain"),
		g("2.  Check sent transaction"),
		g("3.  List local account keys"),
		g("4.  Sending a transaction (add --wait-for-confirm=10 to wait 10 seconds for confirmation)"),
		g("5.  Check a completed transaction receipt"),
		g("6.  Import an account using the mnemonic. Prompts the user to give the mnemonic."),
		g("7.  Import an existing keystore file"),
		g("8.  Import a keystore file using a secp256k1 private key"),
		g("9.  Export a keystore file's secp256k1 private key"),
		g("10. Generate a BLS key then encrypt and save the private key to the specified location. Prompts user to give a password to lock the file."),
		g("11. Create a new validator with a list of BLS keys"),
		g("12. Edit an existing validator"),
		g("13. Delegate an amount to a validator"),
		g("14. Undelegate to a validator"),
		g("15. Collect block rewards as a delegator"),
		g("16. Check active validators"),
		g("17. Check in-memory record of failed staking transactions"),
	)
)
