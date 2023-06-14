package api

import (
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/iyume/dapp-chat/p2pchat/db"
)

const ProtocolName = "p2pchat"

const ProtocolVersion = 1

const protocolLength = 3

const maxMessageSize uint32 = 2 * 1024 * 1024

// uint64 MsgCode
const (
	StatusMsg              = 0x00
	P2PMessageEventMsg     = 0x01
	ChannelMessageEventMsg = 0x02
)

// MakeProtocols always returns latest protocol and drop support for old protocol version.
func MakeProtocols(backend *Backend) []p2p.Protocol {
	protocols := make([]p2p.Protocol, 1)
	protocols[0] = p2p.Protocol{
		Name:    ProtocolName,
		Version: ProtocolVersion,
		Length:  protocolLength,
		Run: func(peer *p2p.Peer, rw p2p.MsgReadWriter) error {
			db.NotateP2PSession(backend.SessionID(peer.ID()), peer.ID(), peer.Node().Pubkey())
			p := NewPeer(peer, rw, ProtocolVersion)
			defer p.Close()
			backend.addPeer(p)
			return Handle(backend, peer, rw)
		},
	}
	return protocols
}
