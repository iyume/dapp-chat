#!/bin/sh
# TODO: rewrite script using python
geth --mine --datadir nodes/signer1 --networkid 12345 --port 30304 --authrpc.port 8552 \
    --unlock 6110a1D3E14FBdD5556F77Edb2785C72D5F50EDb \
    --miner.etherbase 6110a1D3E14FBdD5556F77Edb2785C72D5F50EDb \
    --miner.gasprice 1 --syncmode full \
    --bootnodes 'enode://9b91d8e38cf58aaf6eca6d05627b8d266baabf445e5980194932ed1827c7ff66527aeb4dfe5e62d1e4b4a51f9e5e53c3222966641c0eb9c4dbd1bea946d44d80@127.0.0.1:30303' \
    --nat none --netrestrict '127.0.0.0/8,10.0.0.0/8,172.16.0.0/12,192.168.0.0/16'
