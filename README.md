# Harmony's go-sdk

This is a go layer on top of the Harmony RPC, included is a CLI tool that you can build with a
simple invocation of `make`

See https://docs.harmony.one/home/network/wallets/harmony-cli/download-setup for detailed
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

Cookbook of Usage

Note:

1) Every subcommand recognizes a '--help' flag
2) If a passphrase is used by a subcommand, one can enter their own passphrase interactively
   with the --passphrase option. Alternatively, one can pass their own passphrase via a file
   using the --passphrase-file option. If no passphrase option is selected, the default
   passphrase of '' is used.
3) These examples use Shard 0 of Mainnet as argument for --node

Examples:

1.  Check account balance on given chain
./hmy --node=https://api.s0.t.hmny.io balances <SOME_ONE_ADDRESS>

2.  Check sent transaction
./hmy --node=https://api.s0.t.hmny.io blockchain transaction-by-hash <SOME_TX_HASH>

3.  List local account keys
./hmy keys list

4.  Sending a transaction (waits 40 seconds for transaction confirmation)
./hmy --node=https://api.s0.t.hmny.io transfer \
    --from <SOME_ONE_ADDRESS> --to <SOME_ONE_ADDRESS> \
    --from-shard 0 --to-shard 1 --amount 200 --passphrase

5.  Sending a batch of transactions as dictated from a file (the `--dry-run` options still apply)
./hmy --node=https://api.s0.t.hmny.io transfer --file <PATH_TO_JSON_FILE>
Check README for details on json file format.

6.  Check a completed transaction receipt
./hmy --node=https://api.s0.t.hmny.io blockchain transaction-receipt <SOME_TX_HASH>

7.  Import an account using the mnemonic. Prompts the user to give the mnemonic.
./hmy keys recover-from-mnemonic <ACCOUNT_NAME>

8.  Import an existing keystore file
./hmy keys import-ks <PATH_TO_KEYSTORE_JSON>

9.  Import a keystore file using a secp256k1 private key
./hmy keys import-private-key <secp256k1_PRIVATE_KEY>

10. Export a keystore file's secp256k1 private key
./hmy keys export-private-key <ACCOUNT_ADDRESS> --passphrase

11. Generate a BLS key then encrypt and save the private key to the specified location.
./hmy keys generate-bls-key --bls-file-path <PATH_FOR_BLS_KEY_FILE>

12. Create a new validator with a list of BLS keys
./hmy --node=https://api.s0.t.hmny.io staking create-validator --amount 10 --validator-addr <SOME_ONE_ADDRESS> \
    --bls-pubkeys <BLS_KEY_1>,<BLS_KEY_2>,<BLS_KEY_3> \
    --identity foo --details bar --name baz --max-change-rate 0.1 --max-rate 0.1 --max-total-delegation 10 \
    --min-self-delegation 10 --rate 0.1 --security-contact Leo  --website harmony.one --passphrase

13. Edit an existing validator
./hmy --node=https://api.s0.t.hmny.io staking edit-validator \
    --validator-addr <SOME_ONE_ADDRESS> --identity foo --details bar \
    --name baz --security-contact EK --website harmony.one \
    --min-self-delegation 0 --max-total-delegation 10 --rate 0.1\
    --add-bls-key <SOME_BLS_KEY> --remove-bls-key <OTHER_BLS_KEY> --passphrase

14. Delegate an amount to a validator
./hmy --node=https://api.s0.t.hmny.io staking delegate \
    --delegator-addr <SOME_ONE_ADDRESS> --validator-addr <VALIDATOR_ONE_ADDRESS> \
    --amount 10 --passphrase

15. Undelegate to a validator
./hmy --node=https://api.s0.t.hmny.io staking undelegate \
    --delegator-addr <SOME_ONE_ADDRESS> --validator-addr <VALIDATOR_ONE_ADDRESS> \
    --amount 10 --passphrase

16. Collect block rewards as a delegator
./hmy --node=https://api.s0.t.hmny.io staking collect-rewards \
    --delegator-addr <SOME_ONE_ADDRESS> --passphrase

17. Check elected validators
./hmy --node=https://api.s0.t.hmny.io blockchain validator elected

18. Get current staking utility metrics
./hmy --node=https://api.s0.t.hmny.io blockchain utility-metrics

19. Check in-memory record of failed staking transactions
./hmy --node=https://api.s0.t.hmny.io failures staking

20. Check which shard your BLS public key would be assigned to as a validator
./hmy --node=https://api.s0.t.hmny.io utility shard-for-bls <BLS_PUBLIC_KEY>

21. List Space In Governance
./hmy governance list-space

22. List Proposal In Space Of Governance
./hmy governance list-proposal --space=[space key, example: staking-testnet]

23. View Proposal In Governance
./hmy governance view-proposal --proposal=[proposal hash]

24. New Proposal In Space Of Governance
./hmy governance new-proposal --proposal-yaml=[file path] --key=[key name]
PS: key must first use (hmy keys import-private-key) to import
Yaml example(time is in UTC timezone):
space: staking-testnet
start: 2020-04-16 21:45:12
end: 2020-04-21 21:45:12
choices:
  - yes
  - no
title: this is title
body: |
  this is body,
  you can write mutli line

25. Vote Proposal In Space Of Governance
./hmy governance vote-proposal --proposal=[proposal hash] --choice=[your choice text, eg: yes] --key=[key name]
PS: key must first use (hmy keys import-private-key) to import

26. Enter Console
./hmy command --net=testnet --shard=0
```

# Sending batched transactions

One may find it useful to send a batch of transaction with 1 instance of the binary.
To do this, one can specify a JSON file with the `transaction` subcommand to dictate a batch of transaction to send
off **in sequential order**.

Example:
```
hmy --node="https://api.s1.t.hmny.io/" transfer --file ./batchTransactions.json
```

> Note that the `--wait-for-confirm` and `--dry-run` options still apply when sending batched transactions

## Transfer JSON file format
The JSON file will be a JSON array where each element has the following attributes:

| Key                 | Value-type | Value-description|
| :------------------:|:----------:| :----------------|
| `from`              | string     | [**Required**] Sender's one address, must have key in keystore. |
| `to`                | string     | [**Required**] The receivers one address. |
| `amount`            | string     | [**Required**] The amount to send in $ONE. |
| `from-shard`        | string     | [**Required**] The source shard. |
| `to-shard`          | string     | [**Required**] The destination shard. |
| `passphrase-file`   | string     | [*Optional*] The file path to file containing the passphrase in plain text. If none is provided, check for passphrase string. |
| `passphrase-string` | string     | [*Optional*] The passphrase as a string in plain text. If none is provided, passphrase is ''. |
| `nonce`             | string     | [*Optional*] The nonce of a specific transaction, default uses nonce from blockchain. |
| `gas-price`         | string     | [*Optional*] The gas price to pay in NANO (1e-9 of $ONE), default is 1. |
| `gas-limit`         | string     | [*Optional*] The gas limit, default is 21000. |
| `stop-on-error`     | boolean    | [*Optional*] If true, stop sending transactions if an error occurred, default is false. |
| `true-nonce`        | boolean    | [*Optional*] If true, send transaction using true on-chain nonce. Cannot be used with `nonce`. If none is provided, use tx pool nonce. |

Example of JSON file:

```json
[
  {
    "from": "one103q7qe5t2505lypvltkqtddaef5tzfxwsse4z7",
    "to": "one1zksj3evekayy90xt4psrz8h6j2v3hla4qwz4ur",
    "from-shard" : "0",
    "to-shard": "0",
    "amount": "1",
    "passphrase-string": "",
    "nonce": "35",
    "stop-on-error": true
  },
  {
    "from": "one103q7qe5t2505lypvltkqtddaef5tzfxwsse4z7",
    "to": "one1zksj3evekayy90xt4psrz8h6j2v3hla4qwz4ur",
    "from-shard" : "0",
    "to-shard": "0",
    "amount": "1",
    "passphrase-file": "./pw.txt",
    "gas-price": "1",
    "gas-limit": "21000",
    "true-nonce": true
  }
]
```

## Batched transaction response format

The return will be a JSON array where each element is a transaction log.
The transaction log has the following attributes:

| Key                   | Value-type  | Value-description|
| :--------------------:|:-----------:| :----------------|
| `transaction-receipt` | string      | The transaction hash/receipt if the CLI signed **and sent** a transaction, otherwise this key will not exist |
| `transaction`         | JSON Object | The transaction parameters if `--dry-run` is toggled, otherwise this key will not exist. |
| `blockchain-receipt`  | JSON Object | The transaction receipt from the blockchain if `wait-for-confirm` is > 0, otherwise this key will not exist. |
| `raw-transaction`     | string      | The raw bytes in hex of a sighed transaction if `--dry-run` is toggled, otherwise this key will not exist |
| `errors`              | JSON Array  | A JSON array of strings describing **any** error that occurred during the execution of a transaction. If no errors, this key will not exist. |
| `time-signed-utc`     | string      | The time in UTC as a string of roughly when the transaction was signed. If no signed transaction, this key will not exist. |

Example of returned JSON Array:
```json
[
  {
    "errors": [
      "[2020-01-22 22:01:10.819406] strconv.ParseUint: parsing \"-1\": invalid syntax"
    ]
  },
  {
    "transaction-receipt": "0xf1706080ea9ac210ee2c12c69fb310be5a5da99582b7c783e2f741a3536abbfd",
    "blockchain-receipt": {
      "blockHash": "0xe6de09f4e0ca351257d301a50b4e2ca82646473dc1c4302b570bfab17d421850",
      "blockNumber": "0xb71",
      "contractAddress": null,
      "cumulativeGasUsed": "0x5208",
      "from": "one103q7qe5t2505lypvltkqtddaef5tzfxwsse4z7",
      "gasUsed": "0x5208",
      "logs": [],
      "logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
      "shardID": 0,
      "status": "0x1",
      "to": "one1zksj3evekayy90xt4psrz8h6j2v3hla4qwz4ur",
      "transactionHash": "0xf1706080ea9ac210ee2c12c69fb310be5a5da99582b7c783e2f741a3536abbfd",
      "transactionIndex": "0x0"
    },
    "time-signed-utc": "2020-01-22 22:01:11.468407"
  }
]
```

## Offline sign transfer
1. Get Nonce From a Account. (Need to be online, but no passphrase required)
```bash
./hmy get-nonce --node=https://api.s0.t.hmny.io --from=[ONE address]
```

2. Sign transfer and write to file. (Passphrase required, But no need to be online)
```bash
./hmy transfer --offline-sign --nonce=[nonce value from previous] --from=[ONE address] --to=[ONE address] --amount=1000 --from-shard=0 --to-shard=0 > signed.json
```

3. send `signed.json` to Harmony blockchain! (Need to be online, but no passphrase required)
```bash
./hmy offline-sign-transfer --node=https://api.s0.b.hmny.io --file ./signed.json
```

# Debugging

The go-sdk code respects `HMY_RPC_DEBUG HMY_TX_DEBUG` as debugging
based environment variables.

```bash
HMY_RPC_DEBUG=true HMY_TX_DEBUG=true ./hmy blockchain protocol-version
```

# Contract Deploy

You can deploy the contract use the command;

```bash
./hmy command --node="https://api.s0.b.hmny.io" --net=testnet
```

Then you can test this code to deploy a smart contract

```js
var abi = [{"constant":false,"inputs":[],"name":"changeIt","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"number","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"}]
var code = '0x60806040526212345660005534801561001757600080fd5b5060c8806100266000396000f3006080604052600436106049576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806327b5c17514604e5780638381f58a146062575b600080fd5b348015605957600080fd5b506060608a565b005b348015606d57600080fd5b5060746096565b6040518082815260200191505060405180910390f35b62654321600081905550565b600054815600a165627a7a72305820294db9eed647c48762212fd43bf9ed7733c9d9f43b1ce5583bb4c2bae97fae770029'
var account = 'one1xjanr7lgulc0fqyc8dmfp6jfwuje2d94xfnzyd' // change this to your account address, support eth or one address
var myContract=eth.contract(abi)
contract = myContract.new({from:account,data:code})
```

Wait for a few seconds and it's ready to use