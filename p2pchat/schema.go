package main

// Increase complexity in golang...
type QuerySendP2PMessage struct {
	Nickname string `json:"nickname"`
	UserID   uint64 `json:"user_id"`
	Message  interface{}
}

func (q QuerySendP2PMessage) HasZeroField() bool {
	return q.Nickname == "" || q.UserID == 0 || q.Message == nil
}
