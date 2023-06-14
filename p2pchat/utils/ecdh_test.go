package utils

import (
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

func TestEncryptDecrypt(t *testing.T) {
	assert := assert.New(t)

	key1, err := crypto.GenerateKey()
	assert.Nil(err)
	key2, err := crypto.GenerateKey()
	assert.Nil(err)

	plaintext := []byte("hello")
	ciphertext, err := Encrypt(key1, &key2.PublicKey, plaintext)
	assert.Nil(err)
	ciphertextDec, err := Decrypt(key2, &key1.PublicKey, ciphertext)
	assert.Nil(err)
	assert.Equal(plaintext, ciphertextDec)
}
