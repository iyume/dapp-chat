package main

import (
	"bytes"
	"crypto/ecdsa"
	"log"

	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"github.com/ethereum/go-ethereum/p2p/enode"
)

// source at p2p/discover/v5wire/crypto.go
func ecdh33(privkey *ecdsa.PrivateKey, pubkey *ecdsa.PublicKey) []byte {
	secX, secY := pubkey.Curve.ScalarMult(pubkey.X, pubkey.Y, privkey.D.Bytes())
	log.Printf("secX: %x secY: %x\n", secX.Bytes(), secY.Bytes())
	if secX == nil {
		return nil
	}
	sec := make([]byte, 33)
	sec[0] = 0x02 | byte(secY.Bit(0))
	math.ReadBits(secX, sec[1:])
	return sec
}

func testECDH() {
	// NOTE: crypto/ecdsa PrivateKey.ECDH() create crypto/ecdh PrivateKey
	// but it only supports P256, P384, P521, not S256
	key1, err := crypto.GenerateKey()
	if err != nil {
		log.Fatalln(err)
	}
	key2, err := crypto.GenerateKey()
	if err != nil {
		log.Fatalln(err)
	}
	secret1 := ecdh33(key1, &key2.PublicKey)
	secret2 := ecdh33(key2, &key1.PublicKey)
	log.Printf("secret1: %x\n", secret1)
	log.Printf("secret2: %x\n", secret2)
	log.Println("length:", len(secret1), "equality:", bytes.Equal(secret1, secret2))
}

func testECIES() {
	key1, err := crypto.GenerateKey()
	if err != nil {
		log.Fatalln(err)
	}
	key2, err := crypto.GenerateKey()
	if err != nil {
		log.Fatalln(err)
	}
	key3 := ecies.ImportECDSA(key1)
	key4 := ecies.ImportECDSA(key2)
	keyLen := key3.Params.KeyLen
	log.Println("ECIES S256 sym key length:", keyLen)
	secret1, err := key3.GenerateShared(&key4.PublicKey, keyLen, keyLen)
	if err != nil {
		log.Fatalln(err)
	}
	secret2, err := key4.GenerateShared(&key3.PublicKey, keyLen, keyLen)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("secret1: %x\n", secret1)
	log.Printf("secret2: %x\n", secret2)
	log.Println("length:", len(secret1), "equality:", bytes.Equal(secret1, secret2))
}

func main() {
	testECDH()
	testECIES()
	key, _ := crypto.GenerateKey()
	log.Printf("Keccak256 gen: %x\n", crypto.Keccak256(crypto.FromECDSAPub(&key.PublicKey)[1:]))
	// Convert to [32]byte because enode.ID has String()
	log.Printf("NodeID target: %x\n", enode.PubkeyToIDV4(&key.PublicKey).Bytes())
}
