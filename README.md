# Harmony's go-sdk

This is a go layer on top of the Harmony RPC, included is a CLI tool that you can build with a
simple invocation of `make`

See https://docs.harmony.one/sdk-wiki/command-line-interface/using-the-harmony-cli-tool for detailed
documentation on how to use the `hmy` CLI tools

# Build

Working directly on this repo can be challenging because of the upstream dependencies. Follow the
README in the main repo for an already ready development environment:
https://github.com/harmony-one/harmony/blob/master/README.md.

...for the impatient:

```
$ docker run -it  harmonyone/main:stable /bin/bash
$ cd ../go-sdk
$ git pull -r origin master
$ make
```

# Usage & Examples

`hmy` implements a fluent API, that is, there is a hierarchy of commands.

# bash completions

once built, add `hmy` to your path and add to your `.bashrc`

```
. <(hmy completion)
```

invoke the following command to see the most command usages of `hmy`

```
$ hmy cookbook

Note:

1) Every subcommand recognizes a '--help' flag
2) These examples use shard 1 of testnet as argument for --node

Examples:

1.  Check account balance on given chain
hmy --node="https://api.s1.p.hmny.io/" balances <SOME_ONE_ADDRESS>

2.  Check sent transaction
hmy --node="https://api.s1.p.hmny.io" blockchain transaction-by-hash <SOME_TX_HASH>

3.  List local account keys
hmy keys list

4.  Sending a transaction (add --wait-for-confirm=10 to wait 10 seconds for confirmation)
hmy --node="https://api.s1.p.hmny.io/" transfer \
    --from one1yc06ghr2p8xnl2380kpfayweguuhxdtupkhqzw \
    --to one1q6gkzcap0uruuu8r6sldxuu47pd4ww9w9t7tg6 \
    --from-shard 0 --to-shard 1 --amount 200

5.  Check a completed transaction receipt
hmy --node="https://api.s1.p.hmny.io" blockchain transaction-receipt <SOME_TX_HASH>

6.  Import an account using the mnemonic. Prompts the user to give the mnemonic.
hmy keys add --recover

7.  Import an existing keystore file
hmy keys import-ks <SOME_ABSOLUTE_PATH_TO_KEYSTORE_JSON>.key

8.  Import a keystore file using a secp256k1 private key
hmy keys import-private-key <secp256k1_PRIVATE_KEY>

9.  Export a keystore file's secp256k1 private key
hmy keys export-private-key <ACCOUNT_ADDRESS> --passphrase <YOUR_PASSWORD>

10. Generate a BLS key then encrypt and save the private key to the specified location. Prompts user to give a password to lock the file.
hmy keys generate-bls-key --bls-file-path /tmp/file.key

11. Create a new validator with a list of BLS keys
hmy --node="https://api.s1.p.hmny.io" staking create-validator --amount 10 --validator-addr <SOME_ONE_ADDRESS> \
    --bls-pubkeys <BLS_KEY_1>,<BLS_KEY_2>,<BLS_KEY_3> \
    --identity foo --details bar --name baz --max-change-rate 0.1 --max-rate 0.1 --max-total-delegation 10 \
    --min-self-delegation 10 --rate 0.1 --security-contact Leo  --website harmony.one --passphrase <YOUR_PASSWORD>

12. Edit an existing validator
hmy --node="https://api.s1.p.hmny.io" staking edit-validator \
    --validator-addr <SOME_ONE_ADDRESS> --identity foo --details bar \
    --name baz --security-contact EK --website harmony.one \
    --min-self-delegation 0 --max-total-delegation 10 --rate 0.1\
    --add-bls-key <SOME_BLS_KEY> --remove-bls-key <OTHER_BLS_KEY> --passphrase <YOUR_PASSWORD>

13. Delegate an amount to a validator
hmy --node="https://api.s1.p.hmny.io" staking delegate \
    --delegator-addr <SOME_ONE_ADDRESS> --validator-addr <VALIDATOR_ONE_ADDRESS> \
    --amount 10 --passphrase <YOUR_PASSWORD>

14. Undelegate to a validator
hmy --node="https://api.s1.p.hmny.io" staking undelegate \
    --delegator-addr <SOME_ONE_ADDRESS> --validator-addr <VALIDATOR_ONE_ADDRESS> \
    --amount 10 --passphrase <YOUR_PASSWORD>

15. Collect block rewards as a delegator
hmy --node="https://api.s1.p.hmny.io" staking collect-rewards \
    --delegator-addr <SOME_ONE_ADDRESS> --passphrase <YOUR_PASSWORD>

16. Check active validators
hmy --node="https://api.s1.p.hmny.io" blockchain validator all-active

17. Check in-memory record of failed staking transactions
hmy failures staking

```

# Debugging

The go-sdk code respects `HMY_RPC_DEBUG HMY_TX_DEBUG` as debugging
based environment variables.

```bash
HMY_RPC_DEBUG=true HMY_TX_DEBUG=true ./hmy blockchain protocol-version
```
