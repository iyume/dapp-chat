package server

import (
	"encoding/hex"

	"github.com/iyume/dapp-chat/p2pchat/api"
	"github.com/iyume/dapp-chat/p2pchat/db"
	"github.com/iyume/dapp-chat/p2pchat/utils"
)

// Action protocol-level or internal errors
const (
	CallSuccess int = iota
	CallActionNotFound
	CallArgumentNotEnough
	CallDataInvalid
)

type Caller interface {
	// Call returns error should be handled by server
	Call(string, Getter) (map[string]any, int)
}

// may implemented by codegen

var actions = map[string]func(b *api.Backend, p Getter) (map[string]any, int){
	"get_peers_info": func(b *api.Backend, p Getter) (map[string]any, int) {
		return OK(b.PeersInfo()), 0
	},
	"get_friend_list": func(b *api.Backend, p Getter) (map[string]any, int) {
		resp := []map[string]any{}
		for nodeID, info := range *db.GetFriends() {
			resp = append(resp, map[string]any{
				"node_id": hex.EncodeToString(nodeID[:]),
				"remark":  info.Remark,
			})
		}
		return OK(resp), 0
	},
	"add_friend": func(b *api.Backend, p Getter) (map[string]any, int) {
		nodeIDstr := p.GetString("node_id")
		remark := p.GetString("remark")
		if nodeIDstr == "" || remark == "" {
			return nil, CallArgumentNotEnough
		}
		nodeID, err := utils.ParseNodeID(nodeIDstr)
		if err != nil {
			return Failed("parameter invalid"), CallDataInvalid
		}
		db.AddFriend(nodeID, remark)
		return OK(nil), 0
	},
	"delete_friend": func(b *api.Backend, p Getter) (map[string]any, int) {
		nodeIDstr := p.GetString("node_id")
		if nodeIDstr == "" {
			return nil, CallArgumentNotEnough
		}
		nodeID, err := utils.ParseNodeID(nodeIDstr)
		if err != nil {
			return Failed("parameter invalid"), CallDataInvalid
		}
		db.DeleteFriend(nodeID)
		return OK(nil), 0
	},
	// return json message with any potential error
}

type caller struct {
	backend *api.Backend
}

func (c caller) Call(action string, p Getter) (map[string]any, int) {
	handler, ok := actions[action]
	if !ok {
		return nil, CallActionNotFound
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
