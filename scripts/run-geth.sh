#!/bin/sh
# rm -rf data/geth && geth init --datadir data genesis.json
geth --datadir data --networkid 12345 --port 30303 --authrpc.port 8551 \
    --http --http.addr localhost --http.port 8545 \
    --syncmode full \
    --rpc.txfeecap 0 --rpc.gascap 0 --gpo.ignoreprice 1 \
    --nat none --netrestrict '127.0.0.1/8'
