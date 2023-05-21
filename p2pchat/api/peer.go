package api

import "github.com/ethereum/go-ethereum/p2p"

type Peer struct {
	p       *p2p.Peer
	rw      p2p.MsgReadWriter
	version uint

	closed bool
	term   chan struct{} // received when proto Run is down
}

func NewPeer(p *p2p.Peer, rw p2p.MsgReadWriter, version uint) *Peer {
	return &Peer{
		p:       p,
		rw:      rw,
		version: version,
		term:    make(chan struct{}),
	}
}

func (p *Peer) Close() {
	p.term <- struct{}{}
	p.closed = true
}