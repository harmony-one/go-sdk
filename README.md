# Harmony's go-sdk

This is a go layer on top of the Harmony RPC, included is a CLI tool that you can build with a
simple invocation of `make`

# Build

Invoke `make` to build the `hmy_cli` binary.

# Usage & Examples

hmy_cli implements a fluent API

```bash
./hmy_cli blockchain protocol-version
{"jsonrpc":"2.0","id":"0","result":"0x1"}



./hmy_cli blockchain transaction-by-hash 0x9175e8a3a16afccf6c2d197ed97531de15fcea5b595a4c9ddf4d4d4a22beaab8 --node="http://localhost:9501" --pretty
{
  "id": "0",
  "jsonrpc": "2.0",
  "result": {
    "blockHash": "0x787e5d6570eb7091110239e623a75901fb5d426f2a7ac42b60f83b7b455e200e",
    "blockNumber": "0x100",
    "from": "0xebcd16e8c1d8f493ba04e99a56474122d81a9c58",
    "gas": "0x5208",
    "gasPrice": "0x0",
    "hash": "0x9175e8a3a16afccf6c2d197ed97531de15fcea5b595a4c9ddf4d4d4a22beaab8",
    "input": "0x",
    "nonce": "0x2",
    "r": "0x5bf54b3ae151240d13cfd52acabd78f4161a7f71e5fee203c6ef98719f9ca327",
    "s": "0x129e31ea901f9471597c0c722183ccb0cd7ac628b2831e28068990b087032760",
    "to": "0x514650ca30b3c79f693e14220115434236d44aeb",
    "transactionIndex": "0x0",
    "v": "0x28",
    "value": "0xad78ebc5ac620000"
  }
}

./hmy_cli blockchain transaction-by-receipt 0x8a67a436eadf827faf17bfff18570c1c6c5c63f993d90511364f3adc9dd96c14 --node="http://localhost:9501" --pretty
{
  "id": "0",
  "jsonrpc": "2.0",
  "result": {
    "blockHash": "0xed3b0bb70ab26e5c7db64ad0ccec6c13eef30e2739c5e48b6d676ff33f91e1ac",
    "blockNumber": "0x4c0",
    "contractAddress": null,
    "cumulativeGasUsed": "0x5208",
    "from": "0x514650ca30b3c79f693e14220115434236d44aeb",
    "gasUsed": "0x5208",
    "logs": [],
    "logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
    "shardID": 1,
    "status": "0x1",
    "to": "0x6a87346f3ba9958d08d09484a2b7fdbbe42b0df6",
    "transactionHash": "0x8a67a436eadf827faf17bfff18570c1c6c5c63f993d90511364f3adc9dd96c14",
    "transactionIndex": "0x0"
  }
}
```

# Debugging

```bash
HMY_RPC_DEBUG=true ./hmy_cli blockchain protocol-version
```
