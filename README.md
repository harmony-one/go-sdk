# Harmony's go-sdk

This is a go layer on top of the Harmony RPC, included is a CLI tool that you can build with a
simple invocation of `make`

See https://docs.harmony.one/sdk-wiki/command-line-interface/using-the-harmony-cli-tool for detailed
documentation on how to use the `hmy` CLI tools

# Build

Invoke `make` to build the `hmy` binary.

# Usage & Examples

`hmy` implements a fluent API, that is, there is a hierarchy of commands.

invoke the following command to see the most command usages of `hmy`

```
$ hmy cookbook

Cookbook of usage, note that every subcommand recognizes a '--help' flag

1. Check Balances
hmy --node="https://api.s1.b.hmny.io/" balance <SOME_ONE_ADDRESS>

2. Check completed transaction
hmy --node="https://api.s1.b.hmny.io" blockchain transaction-by-hash <SOME_TRANSACTION_HASH>

3. List local keys
hmy keys list

4. Sending a transaction
./hmy --node="https://api.s0.b.hmny.io/" transfer \
    --from one1yc06ghr2p8xnl2380kpfayweguuhxdtupkhqzw \
    --to one1q6gkzcap0uruuu8r6sldxuu47pd4ww9w9t7tg6 \
    --from-shard 0 --to-shard 1 --amount 200
```

# Debugging

The go-sdk code respects `HMY_RPC_DEBUG HMY_TX_DEBUG` as debugging
based environment variables.

```bash
HMY_RPC_DEBUG=true HMY_TX_DEBUG=true ./hmy blockchain protocol-version
```
