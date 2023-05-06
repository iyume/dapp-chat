#!/bin/sh
geth --datadir data --networkid 12345 --port 30303 --authrpc.port 8551 \
    --http --http.addr 0.0.0.0 --http.port 8545 \
    --syncmode full \
    --rpc.txfeecap 0 --rpc.gascap 0 --gpo.ignoreprice 1 \
    --nat none --netrestrict '127.0.0.1/8'
