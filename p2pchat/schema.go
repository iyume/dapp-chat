package main

type QuerySendP2PMessage struct {
	Nickname string      `json:"nickname"`
	NodeID   string      `json:"user_id"`
	Message  interface{} `json:"message"` // string or json
}

func (q QuerySendP2PMessage) HasZeroField() bool {
	return (q.Nickname == "" && q.NodeID == "") || q.Message == nil || q.Message == ""
}
