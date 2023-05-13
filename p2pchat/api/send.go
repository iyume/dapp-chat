package api

import (
	"log"

	"github.com/ethereum/go-ethereum/p2p"
)

// Backend defines events trigger
type Backend struct {
	running                 bool
	emitP2PMessageEvent     chan *P2PMessageEvent
	emitChannelMessageEvent chan *ChannelMessageEvent
}

// main loop of Backend to be goroutine and call send event methods
func (b *Backend) Run(peer *p2p.Peer, rw p2p.MsgReadWriter) error {
	b.running = true
	for {
		select {
		case p2pevent := <-b.emitP2PMessageEvent:
			p2p.Send(rw, P2PMessageEventMsg, p2pevent)
		case <-b.emitChannelMessageEvent:
			log.Println("send channel message is not supported")
		}
	}
}

func (b *Backend) SendP2PMessageEvent(event *P2PMessageEvent) {
	if !b.running {
		panic("p2p server not started")
	}
	b.emitP2PMessageEvent <- event
}

func (b *Backend) SendChannelMessageEvent(event *ChannelMessageEvent) {
	if !b.running {
		panic("p2p server not started")
	}
	b.emitChannelMessageEvent <- event
}
