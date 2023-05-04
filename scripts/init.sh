#!/bin/sh
rm -rf data/geth nodes/signer1/geth nodes/signer2/geth \
    && geth init --datadir data genesis.json \
    && geth init --datadir nodes/signer1 genesis.json \
    && geth init --datadir nodes/signer2 genesis.json
