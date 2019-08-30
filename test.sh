#!/bin/bash

source ../harmony/scripts/setup_bls_build_flags.sh

# Decent commit is: b4c9a3264a3639367c9baab168aa8e5c7ab2715f
# from harmony repo (needed to check balances, etc)
s='one1tp7xdd9ewwnmyvws96au0e7e7mz6f8hjqr3g3p'
r='one1spshr72utf6rwxseaz339j09ed8p6f8ke370zj'

# --node http://s0.b.hmny.io:9500 \

function check_balances() {
    HMY_RPC_DEBUG=true HMY_TX_DEBUG=true ./hmy_cli account ${s}
    HMY_RPC_DEBUG=true HMY_TX_DEBUG=true ./hmy_cli account ${r}
}

printf '======Balances PRIOR to transfer======\n'
check_balances

HMY_RPC_DEBUG=true HMY_TX_DEBUG=true ./hmy_cli transfer \
	  --from-address=${s} \
	  --to-address=${r} \
	  --from-shard=0 \
	  --to-shard=2 \
	  --amount=10 \
	  --pretty

sleep 5

printf '======Balances AFTER transfer======\n'
check_balances
