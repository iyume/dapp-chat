#!/bin/sh
echo "build constracts..."
solcjs --abi contracts/chat.sol -o build
solcjs --bin contracts/chat.sol -o build
echo "generate go bindings..."
~/go/bin/abigen \
    --abi build/contracts_chat_sol_Chat.abi \
    --bin build/contracts_chat_sol_Chat.bin \
    --pkg chatABI \
    --type Chat \  # same as contract name
    --out go-chat/ABI.go
