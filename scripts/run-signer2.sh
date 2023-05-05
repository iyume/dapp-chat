#!/bin/sh
geth --mine --datadir nodes/signer2 --networkid 12345 --port 30305 --authrpc.port 8553 \
    --unlock 6C4b97f4911951cb79aD44858F8e94D28764d2e5 \
    --miner.etherbase 6C4b97f4911951cb79aD44858F8e94D28764d2e5 \
    --miner.gasprice 1 --syncmode full \
    --bootnodes 'enode://3302e63f80233f1176ceae92dc57df00711a6263b2a6dbbc757121a445963890dd7a821ea143fa3b939b716e1945f091cd3afe54d897a03e02823f26d4b2837e@127.0.0.1:30303' \
    --nat none --netrestrict '127.0.0.1/8'
