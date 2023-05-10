#!/bin/sh
# TODO: rewrite in python script without reset nodekey
rm -rf data/geth/chaindata data/geth/lightchaindata data/geth/transactions.rlp && geth init --datadir data genesis.json
rm -rf nodes/signer1/geth && geth init --datadir nodes/signer1 genesis.json
