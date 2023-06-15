package types

import (
	"time"
)

// 此实现参照 Onebot V12 和 https://github.com/botuniverse/go-libonebots

const (
	EventTypeMessage = "message" // 消息事件
)

// 基础事件类型
type Event struct {
	Time       string `json:"time"`           // RFC3339Nano format
	TimeISO    string `json:"time_iso"`       // RFC3339 format
	Type       string `json:"type"`           // message or other ob12 type
	DetailType string `json:"detail_type"`    // p2p or channel
	Hash       string `json:"hash,omitempty"` // optional hash
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
	UserID string `json:"user_id"` // Sender Node ID, empty in send
}

// MakeP2PMessageEvent do not assiocate with anyone, it should just be sent and handled by receiver
func MakeP2PMessageEvent(message Message) P2PMessageEvent {
	now := time.Now()
	return P2PMessageEvent{
		MessageEvent: MessageEvent{
			Event: Event{
				Time:       now.Format(time.RFC3339Nano),
				TimeISO:    now.Format(time.RFC3339),
				Type:       "message",
				DetailType: "p2p",
			},
			Message: message},
	}
}

type ChannelMessageEvent struct {
	MessageEvent
	ChannelID string `json:"channel_id"` // Channel ID
	UserID    string `json:"user_id"`    // Sender Node ID
}
