package api

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/p2p"
)

// TODO: add friends database, fields: local_user_id, node_id, pubkey

func Handle(backend *Backend, peer *p2p.Peer, rw p2p.MsgReadWriter) error {
	// maybe check protocol-level error
	for {
		msg, err := rw.ReadMsg()
		log.Println("receive Msg")
		if err != nil {
			return err
		}
		if msg.Size > maxMessageSize {
			return fmt.Errorf("message too large: %v > %v", msg.Size, maxMessageSize)
		}
		if err := handleMessage(backend, peer, &msg); err != nil {
			return err
		}
		msg.Discard()
	}
}

var handlerRegistry = map[uint64]func(*p2p.Peer, *p2p.Msg) error{
	StatusMsg:              nil,
	P2PMessageEventMsg:     handleP2PMessage,
	ChannelMessageEventMsg: nil,
}

// handle p2p message, save validated event into database
func handleMessage(backend *Backend, peer *p2p.Peer, msg *p2p.Msg) error {
	// check message is sent to me
	// idslice := make([]byte, 32)
	// if _, err := io.ReadFull(msg.Payload, idslice); err != nil {
	// 	return err
	// }
	// id := [32]byte(idslice)
	// if [32]byte(id) != peer.ID() {
	// 	return errors.New("discard msg not sent to me")
	// }
	// decrypt message and dispatch
	if handler := handlerRegistry[msg.Code]; handler != nil {
		return handler(peer, msg)
	}
	return nil
}

func handleP2PMessage(peer *p2p.Peer, msg *p2p.Msg) error {
	var event P2PMessageEvent
	if err := msg.Decode(&event); err != nil {
		return err
	}
	log.Println("received p2p event:", event)
	// TODO: handle it in database
	return nil
}
