package ipfsutils

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/iyume/dapp-chat/p2pchat/utils"
)

var (
	ErrSelfIDNotMatch = errors.New("~/self_id not matches the given key")
)

func getContent(path string) ([]byte, bool, error) {
	resp, err := http.Get(path)
	if err != nil {
		return nil, false, err
	}
	if resp.StatusCode >= 400 {
		return nil, resp.StatusCode == 404, fmt.Errorf("request %s failed with code %v", path, resp.StatusCode)
	}
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, false, err
	}
	return content, false, nil
}

// ExtractByIndex gains messages in specific p2p session from IPFS gateway. Returns hash to message object.
// Doing this after self_id verified.
func ExtractByIndex(path string, empty_ok bool) (map[string][]byte, error) {
	indexBytes, empty, err := getContent(path + "/index")
	if empty {
		log.Println("index not exists at", path)
		return nil, nil
	}
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println("indexes", path, string(indexBytes))
	res := map[string][]byte{}
	for _, index := range strings.Split(string(indexBytes), "\n") {
		if index == "" {
			continue
		}
		content, _, err := getContent(fmt.Sprintf("%s/%s", path, index))
		if err != nil {
			return nil, err
		}
		res[index] = content
	}
	return res, nil
}

// Verify verifies all messages on self and target IPFS gateway. Returns message hashes.
func Verify(priv *ecdsa.PrivateKey, pubkey *ecdsa.PublicKey, gatewaySelf, gatewayTarget string) ([]string, error) {
	gatewaySelf = strings.TrimSuffix(gatewaySelf, "/")
	gatewayTarget = strings.TrimSuffix(gatewayTarget, "/")

	selfID, _, err := getContent(gatewaySelf + "/self_id")
	if err != nil {
		return nil, err
	}
	targetID, _, err := getContent(gatewayTarget + "/self_id")
	if err != nil {
		return nil, err
	}
	log.Printf("gateway selfID: %s targetID: %s\n", selfID, targetID)
	log.Println("verifying self ID", string(selfID))
	if string(selfID) != enode.PubkeyToIDV4(&priv.PublicKey).String() {
		return nil, ErrSelfIDNotMatch
	}
	log.Println("verifying target ID", string(targetID))
	if string(targetID) != enode.PubkeyToIDV4(pubkey).String() {
		return nil, ErrSelfIDNotMatch
	}

	msgEnc, err := ExtractByIndex(fmt.Sprintf("%s/p2p_sessions/%s", gatewaySelf, targetID), true)
	if err != nil {
		return nil, err
	}
	msgTargetEnc, err := ExtractByIndex(fmt.Sprintf("%s/p2p_sessions/%s", gatewayTarget, selfID), true)
	if err != nil {
		return nil, err
	}
	for key, value := range msgTargetEnc {
		msgEnc[key] = value
	}
	log.Println("message encrypted:", msgEnc)

	var res []string
	for hash, ciphertext := range msgEnc {
		plaintext, err := utils.Decrypt(priv, pubkey, ciphertext)
		if err != nil {
			return nil, err
		}
		hashBytes, err := hex.DecodeString(hash)
		if err != nil {
			return nil, err
		}
		calHash := sha256.Sum256(plaintext)
		if !bytes.Equal(hashBytes, calHash[:]) {
			return nil, err
		} else {
			res = append(res, hash)
		}
	}

	return res, nil
}
