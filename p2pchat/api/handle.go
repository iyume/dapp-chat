package api

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/p2p"
)

func Handle(peer *p2p.Peer, rw p2p.MsgReadWriter) error {
	// TODO: should check protocol-level error
	for {
		msg, err := rw.ReadMsg()
		if err != nil {
			return err
		}
		if msg.Size > maxMessageSize {
			return fmt.Errorf("message too large: %v > %v", msg.Size, maxMessageSize)
		}
		if err := handleMessage(peer, &msg); err != nil {
			return err
		}
		msg.Discard()
	}
}

var handlerRegistry = map[uint64]func(*p2p.Msg) error{
	StatusMsg:              nil,
	P2PMessageEventMsg:     handleP2PMessage,
	ChannelMessageEventMsg: nil,
}

// handle p2p message, save validated event into database
func handleMessage(peer *p2p.Peer, msg *p2p.Msg) error {
	pubkey := peer.Node().Pubkey()
	_ = pubkey
	// decrypt message and dispatch
	if handler := handlerRegistry[msg.Code]; handler != nil {
		return handler(msg)
	}
	return nil
}

func handleP2PMessage(msg *p2p.Msg) error {
	message := P2PMessageEvent{}
	if err := msg.Decode(&message); err != nil {
		return err
	}
	log.Println("receive p2p message:", message)
	// TODO: handle it in database
	return nil
}
