#!/bin/bash

source ../harmony/scripts/setup_bls_build_flags.sh

set -eiu

 # ./hmy.sh transfer --amount 100 --chain-id "testnet" --from one1puj38zamhlu89enzcdjw6rlhlqtyp2c675hjg5 --to one1puj38zamhlu89enzcdjw6rlhlqtyp2c675hjg5 --from-shard 0 --to-shard 1 -n "https://api.s0.b.hmny.io:443"


function test_transfer() {

    local sd='one1yc06ghr2p8xnl2380kpfayweguuhxdtupkhqzw'
    local rcr='one1q6gkzcap0uruuu8r6sldxuu47pd4ww9w9t7tg6'
    local lcl='http://localhost:9500'

    local mainnet='https://api.s0.t.hmny.io/'
    local testnet='https://api.s0.b.hmny.io/'

    local test_net_direct='http://3.80.213.219:9500'
    local test_debug_x1='one1tp7xdd9ewwnmyvws96au0e7e7mz6f8hjqr3g3p'
    local test_debug_x2='one1spshr72utf6rwxseaz339j09ed8p6f8ke370zj'

    local jg5='one1puj38zamhlu89enzcdjw6rlhlqtyp2c675hjg5'

    local spc='https://api.s0.b.hmny.io:443'

    printf 'Before transfer----\nSender:%s Balance\n' ${test_net_direct}
    ./hmy --node=${test_net_direct} balance ${jg5}
    # printf 'Receiver %s Balance:\n' ${rcr}
    # ./hmy --node=${testnet} balance ${test_debug_x2}

    ./hmy --node=${test_net_direct} --verbose \
    	  transfer --from ${jg5} --to ${jg5} \
    	  --from-shard 0 --to-shard 1 --amount 100 --chain-id="mainnet"

    sleep 12

    # printf 'After transfer----\nReceiver:%s Balance\n' ${rcr}
    # ./hmy --node=${beta} balance ${rcr}
    # printf 'Sender %s Balance:\n' ${sd}
    # ./hmy --node=${beta} balance ${sd}

}

test_transfer
