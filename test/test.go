package main

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
)

func main() {
	keyFile := "../data/keystore/UTC--2023-04-18T12-49-23.162529799Z--d542be4551d114a7a2b544bafb7a9feba8784e69"
	password := "123"

	// 读取 keystore 文件
	jsonBytes, err := os.ReadFile(keyFile)
	if err != nil {
		fmt.Println("Error reading keystore file:", err)
		return
	}

	// 解密 keystore 文件
	key, err := keystore.DecryptKey(jsonBytes, password)
	_ = key
	if err != nil {
		fmt.Println("Error decrypting keystore:", err)
		return
	}
}
