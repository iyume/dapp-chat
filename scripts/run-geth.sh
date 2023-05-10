#!/bin/sh
# for firefox, CORS origin should be set
geth --datadir data --networkid 12345 --port 30303 --authrpc.port 8551 \
    --http --http.addr 0.0.0.0 --http.port 8545 \
    --http.corsdomain=moz-extension://d5361047-3a35-4535-af9e-141af00c7fc9 \
    --syncmode full \
    --rpc.txfeecap 0 --rpc.gascap 0 --gpo.ignoreprice 1 \
    --nat none --netrestrict '127.0.0.0/8,10.0.0.0/8,172.16.0.0/12,192.168.0.0/16'
