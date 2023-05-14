package api

import (
	"crypto/ecdsa"
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
	return &Backend{key: key}
}

func (backend *Backend) AddPeer(p *Peer) {
	backend.peers = append(backend.peers, p)
}

func (backend *Backend) findPeer(node_id [32]byte) *Peer {
	// we could check protocol version in further
	for _, p := range backend.peers {
		if p.p.ID() == node_id {
			return p
		}
	}
	return nil
}

// main loop of Backend to be goroutine and call send event methods
func (b *Backend) Run() error {
	b.running = true
	for {
		select {
		case p2pevent := <-b.emitP2PMessageEvent:
			node_id := [32]byte{}
			copy(node_id[:], p2pevent.NodeID)
			peer := b.findPeer(node_id)
			if peer == nil {
				log.Println("error sending p2p message event, peer connection is not established")
				break
			}
			p2p.Send(peer.rw, P2PMessageEventMsg, p2pevent)
		case <-b.emitChannelMessageEvent:
			log.Println("send channel message is not supported")
		}
	}
}

func (b *Backend) SendP2PMessageEvent(node_id string, message Message) error {
	if !b.running {
		return errors.New("p2p server not started")
	}
	if len(node_id) != 64 {
		return errors.New("node_id must be string of length 64")
	}
	event := MakeP2PMessageEvent(time.Now(), message, node_id)
	b.emitP2PMessageEvent <- &event
	return nil
}

func (b *Backend) SendChannelMessageEvent() {
	if !b.running {
		panic("p2p server not started")
	}
}
