package utils

import (
	"encoding/hex"
	"errors"
	"strings"
)

func ConvUint64(val any) (uint64, error) {
	switch val := val.(type) {
	case uint64:
		return val, nil
	case int:
		if val < 0 {
			return 0, errors.New("invalid type")
		}
		return uint64(val), nil
	default:
		return 0, errors.New("invalid type")
	}
}

// Parse hex string trimed "0x" to node ID
func ParseNodeID(in string) ([32]byte, error) {
	return ParseHexNodeID(strings.TrimPrefix(in, "0x"))
}

// Parse hex string without trim "0x" to node ID
func ParseHexNodeID(in string) ([32]byte, error) {
	id := [32]byte{}
	if len(in) != 64 {
		return id, errors.New("want hex string of length 64")
	}
	b, err := hex.DecodeString(in)
	if err != nil {
		return id, err
	}
	copy(id[:], b)
	return id, nil
}
