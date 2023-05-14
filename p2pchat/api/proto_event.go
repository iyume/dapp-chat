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
	Time time.Time `json:"time"`

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

func MakeP2PMessageEvent(time time.Time, message Message, node_id string) P2PMessageEvent {
	return P2PMessageEvent{
		MessageEvent: MessageEvent{Event: Event{
			Time:       time,
			Type:       "message",
			DetailType: "p2p",
		}},
		NodeID: node_id,
	}
}

type ChannelMessageEvent struct {
	MessageEvent
}
