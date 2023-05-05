#!/bin/sh
# TODO: rewrite script using python
geth --mine --datadir nodes/signer1 --networkid 12345 --port 30304 --authrpc.port 8552 \
    --unlock 6110a1D3E14FBdD5556F77Edb2785C72D5F50EDb \
    --miner.etherbase 6110a1D3E14FBdD5556F77Edb2785C72D5F50EDb \
    --miner.gasprice 1 --syncmode full \
    --bootnodes 'enode://3302e63f80233f1176ceae92dc57df00711a6263b2a6dbbc757121a445963890dd7a821ea143fa3b939b716e1945f091cd3afe54d897a03e02823f26d4b2837e@127.0.0.1:30303' \
    --nat none --netrestrict '127.0.0.1/8'
