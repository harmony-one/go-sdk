# Harmony's go-sdk

This is a go layer on top of the Harmony RPC, included is a CLI tool that you can build with a
simple invocation of `make`

# Usage

hmy_cli implements a fluent API

```bash
./hmy_cli blockchain protocol-version
{"jsonrpc":"2.0","id":"0","result":"0x1"}
```

# Debugging

```bash
HMY_RPC_DEBUG=true ./hmy_cli blockchain protocol-version
```
