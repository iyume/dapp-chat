package server

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/iyume/dapp-chat/p2pchat/api"
	"github.com/iyume/dapp-chat/p2pchat/db"
	"github.com/iyume/dapp-chat/p2pchat/types"
	"github.com/iyume/dapp-chat/p2pchat/utils"
)

type Caller interface {
	// Call returns error should be handled by server
	Call(string, Getter) map[string]any
}

// may implemented by codegen

var actions = map[string]func(b *api.Backend, p Getter) map[string]any{
	"get_peers_info": func(b *api.Backend, p Getter) map[string]any {
		return OK(b.PeersInfo())
	},
	"get_friend_list": func(b *api.Backend, p Getter) map[string]any {
		resp := []map[string]any{}
		for nodeID, info := range db.GetFriends() {
			resp = append(resp, map[string]any{
				"node_id": hex.EncodeToString(nodeID[:]),
				"remark":  info.Remark,
			})
		}
		return OK(resp)
	},
	"add_friend": func(b *api.Backend, p Getter) map[string]any {
		if miss := p.Require("node_id", "remark"); miss != "" {
			return Failed(fmt.Sprintf("missing parameter '%s'", miss))
		}
		hexNodeID := p.GetString("node_id")
		remark := p.GetString("remark")
		nodeID, err := utils.ParseNodeID(hexNodeID)
		if err != nil {
			return Failed("parameter invalid")
		}
		db.AddFriend(nodeID, remark)
		return OK(nil)
	},
	"delete_friend": func(b *api.Backend, p Getter) map[string]any {
		if miss := p.Require("node_id"); miss != "" {
			return Failed(fmt.Sprintf("missing parameter '%s'", miss))
		}
		hexNodeID := p.GetString("node_id")
		nodeID, err := utils.ParseNodeID(hexNodeID)
		if err != nil {
			return Failed("parameter invalid")
		}
		db.DeleteFriend(nodeID)
		return OK(nil)
	},
	"get_p2p_session": func(b *api.Backend, p Getter) map[string]any {
		if miss := p.Require("node_id"); miss != "" {
			return Failed(fmt.Sprintf("missing parameter '%s'", miss))
		}
		hexNodeID := p.GetString("node_id")
		nodeID, err := utils.ParseNodeID(hexNodeID)
		if err != nil {
			return Failed("parameter invalid")
		}
		return OK(db.GetP2PSession(b.SessionID(nodeID)))
	},
	"send_p2p_message": func(b *api.Backend, p Getter) map[string]any {
		if miss := p.Require("node_id", "message"); miss != "" {
			return Failed(fmt.Sprintf("missing parameter '%s'", miss))
		}
		hexNodeID := p.GetString("node_id")
		nodeID, err := utils.ParseNodeID(hexNodeID)
		if err != nil {
			return Failed("parameter invalid")
		}
		messageIntf := p.Get("message")
		var message types.Message
		switch val := messageIntf.(type) {
		case string:
			message = types.PlaintextToMessage(val)
		case types.Message:
			message = val
		default:
			return Failed("parameter invalid")
		}
		if message.Empty() {
			return Failed("message is empty")
		}
		if err := b.SendP2PMessage(nodeID, message); err != nil {
			log.Println("error sending p2p message:", err)
			return Failed("internal error")
		}
		return OK(nil)
	},
}

type caller struct {
	backend *api.Backend
}

func (c caller) Call(action string, p Getter) map[string]any {
	handler, ok := actions[action]
	if !ok {
		return nil
	}
	return handler(c.backend, p)
}

func NewCaller(b *api.Backend) Caller {
	return caller{backend: b}
}

func OK(data any) map[string]any {
	return map[string]any{"retcode": 0, "data": data}
}

func Failed(reason string) map[string]any {
	return map[string]any{"retcode": 1, "reason": reason}
}
