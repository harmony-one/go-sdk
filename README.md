# Harmony's go-sdk

This is a go layer on top of the Harmony RPC, included is a CLI tool that you can build with a
simple invocation of `make`

# Build

Invoke `make` to build the `hmy` binary.

# Usage & Examples

`hmy` implements a fluent API, that is, there is a hierarchy of commands.

invoke the following command to see the most command usages of `hmy`

```
$ hmy cookbook
```

# Debugging

The go-sdk code respects `HMY_RPC_DEBUG HMY_TX_DEBUG` as debugging
based environment variables.

```bash
HMY_RPC_DEBUG=true HMY_TX_DEBUG=true ./hmy blockchain protocol-version
```
