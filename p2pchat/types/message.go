package types

type Message []Segment

type Segment struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

type TextSegment struct {
	Text string `json:"text"`
}

func (msg Message) ExtractPlaintext() string {
	var text string
	for _, segment := range msg {
		if segment.Type == "text" {
			data := segment.Data.(TextSegment)
			text += data.Text
		}
	}
	return text
}

func (msg Message) Empty() bool {
	if len(msg) == 0 {
		return true
	}
	if len(msg) == 1 && msg[0].Type == "text" && msg[0].Data.(TextSegment).Text == "" {
		return true
	}
	return false
}

func PlaintextToMessage(text string) Message {
	return Message{Segment{Type: "text", Data: TextSegment{Text: text}}}
}
