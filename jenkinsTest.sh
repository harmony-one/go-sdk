#!/usr/bin/env bash

source ../harmony/scripts/setup_bls_build_flags.sh
make

until $(curl --location --request POST "localhost:9500" \
   --header "Content-Type: application/json" \
   --data '{"jsonrpc":"2.0","method":"net_version","params":[],"id":1}' > HTML_Output)
do
    echo "Trying to connect..."
    sleep 3
done

valid=False
until $valid
do
    result=$(curl --location --request POST "localhost:9500" \
        --header "Content-Type: application/json" \
        --data '{"jsonrpc":"2.0","method":"hmy_getBlockByNumber","params":["0x1", true],"id":1}' \
         | jq | jq 'select(.result)')
    if [ -z "$result" ]; then
        echo "Waiting for localnet to boot..."
        sleep 3
    else
        valid=True
    fi
done

./testHmy.py