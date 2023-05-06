#!/bin/sh
rm -rf data/geth && geth init --datadir data genesis.json
rm -rf nodes/signer1/geth && geth init --datadir nodes/signer1 genesis.json
