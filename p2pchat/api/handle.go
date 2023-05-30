package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/ethereum/go-ethereum/p2p"
	"github.com/iyume/dapp-chat/p2pchat/db"
	"github.com/iyume/dapp-chat/p2pchat/types"
)

// Handle received message, exit with any protocol-level errors
func Handle(backend *Backend, peer *p2p.Peer, rw p2p.MsgReadWriter) error {
	for {
		msg, err := rw.ReadMsg()
		log.Printf("receive msg #%v (%v bytes)\n", msg.Code, msg.Size)
		if err != nil {
			return err
		}
		if msg.Size > maxMessageSize {
			log.Println("message too large, disconnected")
			return fmt.Errorf("message too large: %v (> %v)", msg.Size, maxMessageSize)
		}
		defer msg.Discard()
		// (decrypt) message and dispatch
		if handler := handlerRegistry[msg.Code]; handler != nil {
			if err := handler(backend, peer, msg.Payload); err != nil {
				return err
			}
		} else {
			// unreachable code
			return fmt.Errorf("no handler for msg #%v", msg.Code)
		}
	}
}

var handlerRegistry = map[uint64]func(*Backend, *p2p.Peer, io.Reader) error{
	StatusMsg:              nil,
	P2PMessageEventMsg:     handleP2PMessageEvent,
	ChannelMessageEventMsg: nil,
}

func handleP2PMessageEvent(b *Backend, p *p2p.Peer, r io.Reader) error {
	var event = new(types.P2PMessageEvent)
	if err := json.NewDecoder(r).Decode(event); err != nil {
		return err
	}
	event.UserID = p.ID().String()
	log.Println("received p2p event:", event)
	db.AddP2PMessageEvent(b.SessionID(p.ID()), *event)
	return nil
}
