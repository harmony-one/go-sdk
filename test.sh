#!/bin/bash

source ../harmony/scripts/setup_bls_build_flags.sh

# Decent commit is: b4c9a3264a3639367c9baab168aa8e5c7ab2715f
# from harmony repo (needed to check balances, etc)
s='one1tp7xdd9ewwnmyvws96au0e7e7mz6f8hjqr3g3p'
r='one1a50tun737ulcvwy0yvve0pvu5skq0kjargvhwe'

HMY_RPC_DEBUG=true ./hmy_cli transfer \
	  --from-address=${s} \
	  --to-address=${r} \
	  --from-shard=0 \
	  --to-shard=2 \
	  --amount=10 \
		--pretty
