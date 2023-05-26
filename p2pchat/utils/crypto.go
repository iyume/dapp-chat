package utils

import (
	"crypto/rand"
)

func GenerateToken(len int) []byte {
	hex := make([]byte, len)
	_, err := rand.Reader.Read(hex)
	if err != nil {
		panic(err)
	}
	return hex
}

// func getPrivateKeyJson(file string) *ecdsa.PrivateKey {
// 	keyjson, err := os.ReadFile(file)
// 	if err != nil {
// 		log.Fatalln("Error load keyfile:", err)
// 	}
// 	key, err := keystore.DecryptKey(keyjson, "123")
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	return key.PrivateKey
// }
