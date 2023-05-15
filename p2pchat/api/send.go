package api

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/p2p"
)

// Backend defines events trigger
type Backend struct {
	key                     *ecdsa.PrivateKey
	peers                   []*Peer
	running                 bool
	emitP2PMessageEvent     chan *P2PMessageEvent
	emitChannelMessageEvent chan *ChannelMessageEvent
}

func NewBackend(key *ecdsa.PrivateKey) *Backend {
	return &Backend{
		key:                     key,
		emitP2PMessageEvent:     make(chan *P2PMessageEvent),
		emitChannelMessageEvent: make(chan *ChannelMessageEvent),
	}
}

func (backend *Backend) AddPeer(p *Peer) {
	backend.peers = append(backend.peers, p)
}

func (backend *Backend) findPeer(nodeID [32]byte) *Peer {
	// we could check protocol version in further
	for _, p := range backend.peers {
		pID := p.p.ID()
		if bytes.Equal(pID[:], nodeID[:]) && !p.closed {
			return p
		}
	}
	return nil
}

func truncateNodeID(nodeID [32]byte) string {
	return string(nodeID[:8]) + "..."
}

func stringToIDV4(s string) ([32]byte, error) {
	nodeIDHex, err := hex.DecodeString(s)
	if err != nil {
		return [32]byte{}, err
	}
	nodeID := [32]byte{}
	copy(nodeID[:], nodeIDHex)
	return nodeID, nil
}

// main loop of Backend to be goroutine and call send event methods
func (b *Backend) Run() error {
	b.running = true
	for {
		select {
		case p2pevent := <-b.emitP2PMessageEvent:
			nodeID, err := stringToIDV4(p2pevent.NodeID)
			if err != nil {
				log.Println("node ID is not valid hex string")
				continue
			}
			peer := b.findPeer(nodeID)
			if peer == nil {
				log.Printf(
					"p2p connection to %s is not established\n", truncateNodeID(nodeID))
				continue
			}
			// the peer connection is secure enough, but we could use ECIES/ECDH
			// for futher security
			p2p.Send(peer.rw, P2PMessageEventMsg, p2pevent)
		case <-b.emitChannelMessageEvent:
			log.Println("send channel message is not supported")
		}
	}
}

func (b *Backend) SendP2PMessage(nodeID string, message Message) error {
	if !b.running {
		return errors.New("p2p server not started")
	}
	if len(nodeID) != 64 {
		return errors.New("node ID must be string of length 64")
	}
	event := MakeP2PMessageEvent(time.Now(), message, nodeID)
	log.Println("sending p2p event:", event)
	b.emitP2PMessageEvent <- &event
	return nil
}

func (b *Backend) SendChannelMessage() {
	if !b.running {
		panic("p2p server not started")
	}
}
