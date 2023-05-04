#!/bin/sh
geth --mine --datadir nodes/signer2 --networkid 12345 --port 30305 --authrpc.port 8553 \
    --unlock 6C4b97f4911951cb79aD44858F8e94D28764d2e5 \
    --miner.etherbase 6C4b97f4911951cb79aD44858F8e94D28764d2e5 \
    --miner.gasprice 0 --syncmode full
