#!/bin/sh
# rm -rf data/geth && geth init --datadir data genesis.json
geth --datadir data --networkid 12345
# geth attach data/geth.ipc
