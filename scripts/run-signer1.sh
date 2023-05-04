#!/bin/sh
geth --mine --datadir nodes/signer1 --networkid 12345 --port 30304 --authrpc.port 8552 \
    --unlock 6110a1D3E14FBdD5556F77Edb2785C72D5F50EDb \
    --miner.etherbase 6110a1D3E14FBdD5556F77Edb2785C72D5F50EDb \
    --miner.gasprice 0 --syncmode full
