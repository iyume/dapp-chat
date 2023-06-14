package utils

import (
	"bytes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"log"

	"github.com/emmansun/gmsm/sm4"
)

func RandomBytes(len int) []byte {
	hex := make([]byte, len)
	if _, err := rand.Reader.Read(hex); err != nil {
		panic(err)
	}
	return hex
}

// GetSessionID by sha256(smaller, bigger)
func GetSessionID(id1 [32]byte, id2 [32]byte) [32]byte {
	switch bytes.Compare(id1[:], id2[:]) {
	case 0:
		log.Println("try generating session ID for self, which is potential mistake")
		return [32]byte{}
	case 1:
		id1, id2 = id2, id1
	}
	return sha256.Sum256(append(id1[:], id2[:]...))
}

// generate shared secret of 32 byte
func ecdhGenerateShared(priv *ecdsa.PrivateKey, pubkey *ecdsa.PublicKey) ([]byte, error) {
	secX, _ := pubkey.ScalarMult(pubkey.X, pubkey.Y, priv.D.Bytes())
	if secX == nil {
		return nil, errors.New("shared point at infinite")
	}
	hash := sha256.Sum256(secX.Bytes())
	return hash[:], nil
}

// Encrypt gets a hashed shared point X as secret, then sym encrypts the plaintext using sm4 (CTR mode)
// The resulting ciphertext length is BlockSize+len(plaintext), because IV is randomly generated and included at the beginning of ciphertext
func Encrypt(priv *ecdsa.PrivateKey, pubkey *ecdsa.PublicKey, plaintext []byte) ([]byte, error) {
	secret, err := ecdhGenerateShared(priv, pubkey)
	if err != nil {
		return nil, err
	}
	sm4Secret := secret[:sm4.BlockSize]

	block, err := sm4.NewCipher(sm4Secret)
	if err != nil {
		panic(err)
	}

	ciphertext := make([]byte, len(plaintext))
	iv := RandomBytes(sm4.BlockSize)

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext, plaintext)

	// Include IV at the beginning of ciphertext
	return append(iv, ciphertext...), nil
}

// Decrypt gets a hashed shared point X as secret, then sym decrypts the ciphertext using sm4 (CTR mode)
func Decrypt(priv *ecdsa.PrivateKey, pubkey *ecdsa.PublicKey, ciphertext []byte) ([]byte, error) {
	if len(ciphertext) <= sm4.BlockSize {
		return nil, fmt.Errorf("length of ciphertext must be more than %v", sm4.BlockSize)
	}

	secret, err := ecdhGenerateShared(priv, pubkey)
	if err != nil {
		return nil, err
	}
	sm4Secret := secret[:sm4.BlockSize]

	block, err := sm4.NewCipher(sm4Secret)
	if err != nil {
		panic(err)
	}

	plaintext := make([]byte, len(ciphertext)-sm4.BlockSize)
	iv := ciphertext[:sm4.BlockSize]

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(plaintext, ciphertext[sm4.BlockSize:])

	return plaintext, nil
}
