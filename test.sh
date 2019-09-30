#!/bin/bash

source ../harmony/scripts/setup_bls_build_flags.sh

set -eiu

function test_transfer() {

    local sd='one1yc06ghr2p8xnl2380kpfayweguuhxdtupkhqzw'
    local rcr='one1q6gkzcap0uruuu8r6sldxuu47pd4ww9w9t7tg6'
    local lcl='http://localhost:9500'
    local main='https://api.s0.t.hmny.io/'
    local beta='https://api.s0.b.hmny.io/'
    local mainDirect='http://18.237.68.133:9500'
    local test_debug_x1='one1tp7xdd9ewwnmyvws96au0e7e7mz6f8hjqr3g3p'
    local test_debug_x2='one1spshr72utf6rwxseaz339j09ed8p6f8ke370zj'

    printf 'Before transfer----\nSender:%s Balance\n' ${sd}
    ./hmy --node=${lcl} balance ${test_debug_x1}
    printf 'Receiver %s Balance:\n' ${rcr}
    ./hmy --node=${lcl} balance ${test_debug_x2}

    HMY_RPC_DEBUG=true HMY_TX_DEBUG=true ./hmy --node=${lcl} \
    	  transfer --from ${test_debug_x2} --to ${test_debug_x1} --chain-id='mainnet' \
    	  --from-shard 0 --to-shard 0 --amount 1

    sleep 12

    printf 'After transfer----\nReceiver:%s Balance\n' ${rcr}
    ./hmy --node=${main} balance ${rcr}
    printf 'Sender %s Balance:\n' ${sd}
    ./hmy --node=${main} balance ${sd}

}

test_transfer
