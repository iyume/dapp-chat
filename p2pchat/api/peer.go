package api

import (
	"bytes"
	"encoding/json"

	"github.com/ethereum/go-ethereum/p2p"
	"github.com/iyume/dapp-chat/p2pchat/types"
)

type Peer struct {
	p       *p2p.Peer
	rw      p2p.MsgReadWriter
	version uint

	closed bool
	term   chan struct{} // received when proto Run is down

	// TODO: message queue
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
	// add this on message queue
	// p.term <- struct{}{}
	p.closed = true
}

// msgpack support?

func SendJson(w p2p.MsgWriter, msgcode uint64, data interface{}) error {
	payload := bytes.NewBuffer(nil)
	if err := json.NewEncoder(payload).Encode(data); err != nil {
		return err
	}
	return w.WriteMsg(p2p.Msg{Code: msgcode, Size: uint32(payload.Len()), Payload: payload})
}

func (p *Peer) SendMessage(event types.P2PMessageEvent) error {
	return SendJson(p.rw, P2PMessageEventMsg, event)
}
