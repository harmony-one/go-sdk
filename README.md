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

Cookbook of Usage

Note:

1) Every subcommand recognizes a '--help' flag
2) If a passphrase is used by a subcommand, one can enter their own passphrase interactively
   with the --passphrase option. Alternatively, one can pass their own passphrase via a file
   using the --passphrase-file option. If no passphrase option is selected, the default
   passphrase of '' is used.
3) These examples use shard 1 of testnet as argument for --node

Examples:

1.  Check account balance on given chain
hmy --node="https://api.s1.t.hmny.io/" balances <SOME_ONE_ADDRESS>

2.  Check sent transaction
hmy --node="https://api.s1.t.hmny.io" blockchain transaction-by-hash <SOME_TX_HASH>

3.  List local account keys
hmy keys list

4.  Sending a transaction (waits 40 seconds for transaction confirmation)
hmy --node="https://api.s1.t.hmny.io/" transfer \
    --from one1yc06ghr2p8xnl2380kpfayweguuhxdtupkhqzw \
    --to one1q6gkzcap0uruuu8r6sldxuu47pd4ww9w9t7tg6 \
    --from-shard 0 --to-shard 1 --amount 200

5.  Sending a batch of transactions as dictated from a file (the `--dry-run` options still apply)
hmy --node="https://api.s1.t.hmny.io/" transfer --file <PATH_TO_JSON_FILE>

    Example of JSON file format:
      [
        {
          "from": "one103q7qe5t2505lypvltkqtddaef5tzfxwsse4z7",
          "to": "one1zksj3evekayy90xt4psrz8h6j2v3hla4qwz4ur",
          "from-shard" : "0",
          "to-shard": "0",
          "amount": "1",
          "passphrase-string": "",
          "nonce": "1",
          "stop-on-error": true
        },
        {
          "from": "one103q7qe5t2505lypvltkqtddaef5tzfxwsse4z7",
          "to": "one1zksj3evekayy90xt4psrz8h6j2v3hla4qwz4ur",
          "from-shard" : "0",
          "to-shard": "0",
          "amount": "1",
          "passphrase-file": "./pw.txt"
        }
      ]

6.  Check a completed transaction receipt
hmy --node="https://api.s1.t.hmny.io" blockchain transaction-receipt <SOME_TX_HASH>

7.  Import an account using the mnemonic. Prompts the user to give the mnemonic.
hmy keys recover-from-mnemonic <ACCOUNT_NAME>

8.  Import an existing keystore file
hmy keys import-ks <PATH_TO_KEYSTORE_JSON>.key

9.  Import a keystore file using a secp256k1 private key
hmy keys import-private-key <secp256k1_PRIVATE_KEY>

10.  Export a keystore file's secp256k1 private key
hmy keys export-private-key <ACCOUNT_ADDRESS> --passphrase

11. Generate a BLS key then encrypt and save the private key to the specified location.
hmy keys generate-bls-key --bls-file-path /tmp/file.key

12. Create a new validator with a list of BLS keys
hmy --node="https://api.s0.t.hmny.io" staking create-validator --amount 10 --validator-addr <SOME_ONE_ADDRESS> \
    --bls-pubkeys <BLS_KEY_1>,<BLS_KEY_2>,<BLS_KEY_3> \
    --identity foo --details bar --name baz --max-change-rate 0.1 --max-rate 0.1 --max-total-delegation 10 \
    --min-self-delegation 10 --rate 0.1 --security-contact Leo  --website harmony.one --passphrase

13. Edit an existing validator
hmy --node="https://api.s0.t.hmny.io" staking edit-validator \
    --validator-addr <SOME_ONE_ADDRESS> --identity foo --details bar \
    --name baz --security-contact EK --website harmony.one \
    --min-self-delegation 0 --max-total-delegation 10 --rate 0.1\
    --add-bls-key <SOME_BLS_KEY> --remove-bls-key <OTHER_BLS_KEY> --passphrase

14. Delegate an amount to a validator
hmy --node="https://api.s0.t.hmny.io" staking delegate \
    --delegator-addr <SOME_ONE_ADDRESS> --validator-addr <VALIDATOR_ONE_ADDRESS> \
    --amount 10 --passphrase

15. Undelegate to a validator
hmy --node="https://api.s0.t.hmny.io" staking undelegate \
    --delegator-addr <SOME_ONE_ADDRESS> --validator-addr <VALIDATOR_ONE_ADDRESS> \
    --amount 10 --passphrase

16. Collect block rewards as a delegator
hmy --node="https://api.s0.t.hmny.io" staking collect-rewards \
    --delegator-addr <SOME_ONE_ADDRESS> --passphrase

17. Check active validators
hmy --node="https://api.s0.t.hmny.io" blockchain validator all-active

18. Get current staking utility metrics
hmy --node="https://api.s0.t.hmny.io" blockchain utility-metrics

19. Check in-memory record of failed staking transactions
hmy failures staking

20. Check which shard your BLS public key would be assigned to as a validator
hmy utility shard-for-bls 2d61379e44a772e5757e27ee2b3874254f56073e6bd226eb8b160371cc3c18b8c4977bd3dcb71fd57dc62bf0e143fd08
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
    "gas-limit": "21000"
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

# Debugging

The go-sdk code respects `HMY_RPC_DEBUG HMY_TX_DEBUG` as debugging
based environment variables.

```bash
HMY_RPC_DEBUG=true HMY_TX_DEBUG=true ./hmy blockchain protocol-version
```
