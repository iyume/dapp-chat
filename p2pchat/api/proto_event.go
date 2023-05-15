package api

import (
	"time"
)

// 此实现参照 Onebot V12 和 https://github.com/botuniverse/go-libonebot
// 去除了很多冗余字段如 user_id, self_id, message_id

const (
	EventTypeMessage = "message" // 消息事件
)

// 基础事件类型
type Event struct {
	// ID   string    `json:"id"` // fake field
	Time string `json:"time"` // the rlp encoding cannot handle time properly

	Type string `json:"type"`

	// p2p or channel
	DetailType string `json:"detail_type"`
}

func (e *Event) Name() string {
	return e.Type + "." + e.DetailType
}

type MessageEvent struct {
	Event
	Message Message `json:"message"`
}

type P2PMessageEvent struct {
	MessageEvent
	NodeID string `json:"node_id"`
}

func MakeP2PMessageEvent(time_ time.Time, message Message, nodeID string) P2PMessageEvent {
	return P2PMessageEvent{
		MessageEvent: MessageEvent{
			Event: Event{
				Time:       time_.Format(time.RFC3339Nano),
				Type:       "message",
				DetailType: "p2p",
			},
			Message: message},
		NodeID: nodeID,
	}
}

type ChannelMessageEvent struct {
	MessageEvent
}
