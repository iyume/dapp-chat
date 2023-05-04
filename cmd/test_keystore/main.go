package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
)

func testDecryptKey(keyfile string, passphrase string) {
	// 读取 keystore 文件 (json bytes)
	keyjson, err := os.ReadFile(keyfile)
	if err != nil {
		log.Fatalln("Error reading keystore file:", err)
	}
	// 解密 keystore 文件
	key, err := keystore.DecryptKey(keyjson, passphrase)
	if err != nil {
		log.Fatalln("Error decrypting keystore:", err)
	}
	fmt.Println("Decrypted key:", key)
}

func getKeyStore(keystoredir string, lightKDF bool) *keystore.KeyStore {
	n, p := keystore.StandardScryptN, keystore.StandardScryptP
	if lightKDF {
		n, p = keystore.LightScryptN, keystore.LightScryptP
	}
	ks := keystore.NewKeyStore(keystoredir, n, p)
	return ks
}

func getKeyStoreStandardKDF(keystoredir string) *keystore.KeyStore {
	return getKeyStore(keystoredir, false)
}

func testFindByAddress(ks *keystore.KeyStore, address common.Address) {
	// find by URL is also OK
	account, err := ks.Find(accounts.Account{Address: address})
	if err != nil {
		log.Fatalln("Cannot find account:", err)
	}
	fmt.Printf("Keyfile for address %v is %v", address, account.URL)
}

func main() {
	testDecryptKey("data/keystore/UTC--2023-04-18T12-49-23.162529799Z--d542be4551d114a7a2b544bafb7a9feba8784e69", "123")
	ks := getKeyStoreStandardKDF("data/keystore")
	fmt.Println("List all accounts:", ks.Accounts())
	testFindByAddress(ks, common.HexToAddress("8cedb7c6af8a7781ec89bb84900768de99c8235b"))
}
