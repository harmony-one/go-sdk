#!/bin/bash

source ../harmony/scripts/setup_bls_build_flags.sh

function test_transfer() {

    local sd='one1yc06ghr2p8xnl2380kpfayweguuhxdtupkhqzw'
    local rcr='one1q6gkzcap0uruuu8r6sldxuu47pd4ww9w9t7tg6'
    local lcl='http://localhost:9500'
    local main='https://api.s0.t.hmny.io/'
    local beta='https://api.s0.b.hmny.io/'
    local mainDirect='http://18.237.68.133:9500'

    printf 'Before transfer----\nSender:%s Balance\n' ${sd}
    ./hmy --node=${main} balance ${sd}
    printf 'Receiver %s Balance:\n' ${rcr}
    ./hmy --node=${main} balance ${rcr}

    HMY_RPC_DEBUG=true HMY_TX_DEBUG=true ./hmy --node=${main} \
    	  transfer --from ${sd} --to ${rcr} \
    	  --from-shard 0 --to-shard 0 --amount 1 --passphrase=''

    sleep 12

    printf 'After transfer----\nReceiver:%s Balance\n' ${rcr}
    ./hmy --node=${main} balance ${rcr}
    printf 'Sender %s Balance:\n' ${sd}
    ./hmy --node=${main} balance ${sd}

}

test_transfer
