package api

import (
	"time"
)

// 此实现参照 Onebot V12 和 https://github.com/botuniverse/go-libonebot
// 所有 ID 相关的数据都不会在发送时构造，而是在接收时构造

const (
	EventTypeMessage = "message" // 消息事件
)

// 基础事件类型
type Event struct {
	Time       string `json:"time"`
	Type       string `json:"type"`        // message or other ob12 type
	DetailType string `json:"detail_type"` // p2p or channel
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
	NodeID string `json:"node_id"` // Sender Node ID
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
	ChannelID string `json:"channel_id"` // Channel ID
	NodeID    string `json:"node_id"`    // Sender Node ID
}
