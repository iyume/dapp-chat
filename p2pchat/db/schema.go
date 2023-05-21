package db

import (
	"encoding/binary"
)

var (
	// localKey register immutable local node ID
	localKey = []byte("LocalID")

	// "F-"
	friendPrefix = []byte("F-")

	// "M-"
	p2pMessageLookupPrefix = []byte("M-")

	// "S-"
	p2pSessionPrefix = []byte("S-")
)

// append encoded uint64 in big endian
func appendUint64(prefix []byte, number uint64) []byte {
	return binary.BigEndian.AppendUint64(prefix, number)
}

// "F-" + [32]byte -> friend stock
func friendKey(nodeID [32]byte) []byte {
	return append(friendPrefix, nodeID[:]...)
}

// "M-" + [8]byte -> session ID
func p2pmessageLookupKey(messageID uint64) []byte {
	return appendUint64(p2pMessageLookupPrefix, messageID)
}

// "S-" + [32]byte -> session stock
func p2pSessionKey(sessionID [32]byte) []byte {
	return append(p2pSessionPrefix, sessionID[:]...)
}
