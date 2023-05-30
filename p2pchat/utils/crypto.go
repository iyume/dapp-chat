package utils

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
)

func GenerateToken(len int) []byte {
	hex := make([]byte, len)
	_, err := rand.Reader.Read(hex)
	if err != nil {
		panic(err)
	}
	return hex
}

// GetSessionID by sha256(smaller, bigger)
func GetSessionID(id1 [32]byte, id2 [32]byte) [32]byte {
	switch bytes.Compare(id1[:], id2[:]) {
	case 0:
		return [32]byte{}
	case 1:
		id1, id2 = id2, id1
	}
	return sha256.Sum256(append(id1[:], id2[:]...))
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
